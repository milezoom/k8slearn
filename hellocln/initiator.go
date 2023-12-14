package hellocln

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"hellocln/contract"
	"hellocln/implementor"
	"hellocln/mocker"

	"git.bluebird.id/bbd/lib/iresolver"
	commonGrpc "git.bluebird.id/mybb-ms/aphrodite/grpc"
	commonMicro "git.bluebird.id/mybb-ms/aphrodite/microservice"
	"go.elastic.co/apm/module/apmgrpc/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func NewMocker() (*mocker.Mocker, error) {
	return &mocker.Mocker{}, nil
}

type ImplementorConfig struct {
	Host     string
	Port     int
	CertPath string
}

func NewImplementor(conf ImplementorConfig) (*implementor.HelloServiceImplementor, error) {
	baseUrl := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	opts := []grpc.DialOption{}

	if conf.CertPath != "" {
		cred, err := commonGrpc.TLSCredentialFromCertForClient(conf.CertPath)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}

	if conf.CertPath == "" && conf.Port == 443 {
		opts = append(opts, grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: true,
			}),
		))
	}

	if conf.CertPath == "" && conf.Port != 443 {
		opts = append(opts, grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		))
	}

	cb := commonMicro.NewCircuitBreaker(commonMicro.DefaultBreakerSetting("hellocln-breaker", 30*time.Second))
	opts = append(opts, grpc.WithUnaryInterceptor(commonGrpc.BreakerClientUnaryInterceptor(cb)))
	opts = append(opts, grpc.WithUnaryInterceptor(apmgrpc.NewUnaryClientInterceptor()))
	opts = append(opts, grpc.WithStreamInterceptor(apmgrpc.NewStreamClientInterceptor()))

	if strings.Contains(conf.Host, "kubernetes://") {
		iresolver.RegisterInCluster()
		opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	}

	conn, err := grpc.Dial(baseUrl, opts...)
	if err != nil {
		return nil, err
	}
	return &implementor.HelloServiceImplementor{
		GrpcConn: conn,
		Cli:      contract.NewHelloServiceClient(conn),
	}, nil
}
