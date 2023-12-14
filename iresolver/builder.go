/*
 *
 * Copyright 2023 author.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author : Robbi Ibadi
 */

package iresolver

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
)

const (
	kubernetesSchema = "kubernetes"
	defaultFreq      = time.Minute * 30
)

type targetInfo struct {
	serviceName       string
	serviceNamespace  string
	port              string
	resolveByPortName bool
	useFirstPort      bool
}

func (ti targetInfo) String() string {
	return fmt.Sprintf("kubernetes://%s/%s:%s", ti.serviceNamespace, ti.serviceName, ti.port)
}

// RegisterInCluster registers the iresolver builder to grpc with kubernetes schema
func RegisterInCluster() {
	RegisterInClusterWithSchema(kubernetesSchema)
}

// RegisterInClusterWithSchema registers the iresolver builder to the grpc with custom schema
func RegisterInClusterWithSchema(schema string) {
	resolver.Register(NewBuilder(nil, schema))
}

// NewBuilder creates a kubeBuilder which is used by grpc resolver.
func NewBuilder(client K8sClient, schema string) resolver.Builder {
	return &kubeBuilder{
		k8sClient: client,
		schema:    schema,
	}
}

type kubeBuilder struct {
	k8sClient K8sClient
	schema    string
}

func splitServicePortNamespace(hpn string) (service, port, namespace string) {
	service = hpn

	colon := strings.LastIndexByte(service, ':')
	if colon != -1 {
		service, port = service[:colon], service[colon+1:]
	}

	parts := strings.SplitN(service, ".", 3)
	if len(parts) >= 2 {
		service, namespace = parts[0], parts[1]
	}

	return
}

func parseResolverTarget(target resolver.Target) (targetInfo, error) {
	var service, port, namespace string

	url, err := parseTargetUrl(formatTarget(target))
	if err != nil {
		return targetInfo{}, fmt.Errorf("target %s must specify a service", &url)
	}
	endpoint := getEndpoint(url)

	if url.Host == "" {
		// kubernetes:///service.namespace:port
		service, port, namespace = splitServicePortNamespace(endpoint)
	} else if url.Port() == "" && endpoint != "" {
		// kubernetes://namespace/service:port
		service, port, _ = splitServicePortNamespace(endpoint)
		namespace = url.Hostname()
	} else {
		// kubernetes://service.namespace:port
		service, port, namespace = splitServicePortNamespace(url.Host)
	}

	if service == "" {
		return targetInfo{}, fmt.Errorf("target %s must specify a service", &url)
	}

	resolveByPortName := false
	useFirstPort := false
	if port == "" {
		useFirstPort = true
	} else if _, err := strconv.Atoi(port); err != nil {
		resolveByPortName = true
	}

	return targetInfo{
		serviceName:       service,
		serviceNamespace:  namespace,
		port:              port,
		resolveByPortName: resolveByPortName,
		useFirstPort:      useFirstPort,
	}, nil
}

func formatTarget(target resolver.Target) string {
	// Reconstruct the URL with the provided resolver.Target fields.
	var endpoint string
	if target.Endpoint() != "" {
		endpoint = "/" + target.Endpoint()
	}

	urlStr := fmt.Sprintf("%s://%s%s", target.URL.Scheme, target.URL.Host, endpoint)

	return urlStr
}

func getEndpoint(url url.URL) string {
	endpoint := url.Path
	if endpoint == "" {
		endpoint = url.Opaque
	}
	return strings.TrimPrefix(endpoint, "/")
}

func parseTargetUrl(target string) (url.URL, error) {
	u, err := url.Parse(target)
	if err != nil {
		return url.URL{}, err
	}
	return *u, nil
}

func (b *kubeBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	if b.k8sClient == nil {
		if cl, err := NewInClusterK8sClient(); err == nil {
			b.k8sClient = cl
		} else {
			return nil, err
		}
	}
	ti, err := parseResolverTarget(target)
	if err != nil {
		return nil, err
	}
	if ti.serviceNamespace == "" {
		ti.serviceNamespace = getCurrentNamespaceOrDefault()
	}
	ctx, cancel := context.WithCancel(context.Background())
	r := &kResolver{
		target:    ti,
		ctx:       ctx,
		cancel:    cancel,
		cc:        cc,
		rn:        make(chan struct{}, 1),
		k8sClient: b.k8sClient,
		t:         time.NewTimer(defaultFreq),
		freq:      defaultFreq,
	}
	go until(func() {
		r.wg.Add(1)
		err := r.watch()
		if err != nil && err != io.EOF {
			grpclog.Errorf("iresolver: watching ended with error='%v', will reconnect again", err)
		}
	}, time.Second, time.Second*30, ctx.Done())
	return r, nil
}

func (b *kubeBuilder) Scheme() string {
	return b.schema
}

type kResolver struct {
	target    targetInfo
	ctx       context.Context
	cancel    context.CancelFunc
	cc        resolver.ClientConn
	rn        chan struct{}
	k8sClient K8sClient
	wg        sync.WaitGroup
	t         *time.Timer
	freq      time.Duration
}

func (k *kResolver) ResolveNow(resolver.ResolveNowOptions) {
	select {
	case k.rn <- struct{}{}:
	default:
	}
}

func (k *kResolver) Close() {
	k.cancel()
	k.wg.Wait()
}

func (k *kResolver) makeAddresses(e Endpoints) ([]resolver.Address, string) {
	var newAddrs []resolver.Address
	for _, subset := range e.Subsets {
		port := ""
		if k.target.useFirstPort {
			port = strconv.Itoa(subset.Ports[0].Port)
		} else if k.target.resolveByPortName {
			for _, p := range subset.Ports {
				if p.Name == k.target.port {
					port = strconv.Itoa(p.Port)
					break
				}
			}
		} else {
			port = k.target.port
		}

		if len(port) == 0 {
			port = strconv.Itoa(subset.Ports[0].Port)
		}

		for _, address := range subset.Addresses {
			newAddrs = append(newAddrs, resolver.Address{
				Addr:       net.JoinHostPort(address.IP, port),
				ServerName: fmt.Sprintf("%s.%s", k.target.serviceName, k.target.serviceNamespace),
				Metadata:   nil,
			})
		}
	}
	return newAddrs, ""
}

func (k *kResolver) handle(e Endpoints) {
	result, _ := k.makeAddresses(e)
	if len(result) > 0 {
		k.cc.NewAddress(result)
	}
}

func (k *kResolver) resolve() {
	e, err := getEndpoints(k.k8sClient, k.target.serviceNamespace, k.target.serviceName)
	if err == nil {
		k.handle(e)
	} else {
		grpclog.Errorf("iresolver: lookup endpoints failed: %v", err)
	}
	k.t.Reset(k.freq)
}

func (k *kResolver) watch() error {
	defer k.wg.Done()
	sw, err := watchEndpoints(k.ctx, k.k8sClient, k.target.serviceNamespace, k.target.serviceName)
	if err != nil {
		return err
	}
	for {
		select {
		case <-k.ctx.Done():
			return nil
		case <-k.t.C:
			k.resolve()
		case <-k.rn:
		case up, hasMore := <-sw.ResultChan():
			if hasMore {
				k.handle(up.Object)
			} else {
				return nil
			}
		}
	}
}
