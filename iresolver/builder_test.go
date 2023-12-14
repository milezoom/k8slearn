package iresolver

import (
	"net/url"
	"strings"
	"testing"

	"google.golang.org/grpc/resolver"
)

func TestParseResolverTarget1(t *testing.T) {
	tests := []struct {
		name   string
		target resolver.Target
		want   targetInfo
		err    bool
	}{
		// Test cases go here
		{
			name: "Test case 1",
			target: resolver.Target{
				URL: url.URL{
					Scheme: "kubernetes",
					Path:   "grpc-server.default:50",
				},
			},
			want: targetInfo{
				serviceName:       "grpc-server",
				serviceNamespace:  "default",
				port:              "50",
				resolveByPortName: false,
				useFirstPort:      false,
			},
			err: false,
		},
		{
			name: "Test case 2",
			target: resolver.Target{
				URL: url.URL{
					Scheme: "kubernetes",
					Path:   "grpc-greeter-server:50051",
				},
			},
			want: targetInfo{
				serviceName:       "grpc-greeter-server",
				serviceNamespace:  "",
				port:              "50051",
				resolveByPortName: false,
				useFirstPort:      false,
			},
			err: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseResolverTarget(tt.target)
			if err == nil && tt.err {
				t.Errorf("case %s: want error but got nil", tt.name)
				// continue
			}
			if err != nil && !tt.err {
				t.Errorf("case %s: got '%v' error but don't want an error", tt.name, err)
				// continue
			}
			if got != tt.want {
				t.Errorf("case %s: parseResolverTarget(%q) = %+v, want %+v", tt.name, tt.target.Endpoint(), got, tt.want)
			}
		})
	}
}

func TestParseTargets(t *testing.T) {
	for i, test := range []struct {
		target string
		want   targetInfo
		err    bool
	}{
		{"", targetInfo{}, true},
		{"kubernetes:///", targetInfo{}, true},
		{"kubernetes://a:30", targetInfo{"a", "", "30", false, false}, false},
		{"kubernetes://a/", targetInfo{"a", "", "", false, true}, false},
		{"kubernetes:///a", targetInfo{"a", "", "", false, true}, false},
		{"kubernetes://a/b", targetInfo{"b", "a", "", false, true}, false},
		{"kubernetes://a.b/", targetInfo{"a", "b", "", false, true}, false},
		{"kubernetes:///a.b:80", targetInfo{"a", "b", "80", false, false}, false},
		{"kubernetes:///a.b:port", targetInfo{"a", "b", "port", true, false}, false},
		{"kubernetes:///a:port", targetInfo{"a", "", "port", true, false}, false},
		{"kubernetes://x/a:port", targetInfo{"a", "x", "port", true, false}, false},
		{"kubernetes://a.x:30/", targetInfo{"a", "x", "30", false, false}, false},
		{"kubernetes://a.b.svc.cluster.local", targetInfo{"a", "b", "", false, true}, false},
		{"kubernetes://a.b.svc.cluster.local:80", targetInfo{"a", "b", "80", false, false}, false},
		{"kubernetes:///a.b.svc.cluster.local", targetInfo{"a", "b", "", false, true}, false},
		{"kubernetes:///a.b.svc.cluster.local:80", targetInfo{"a", "b", "80", false, false}, false},
		{"kubernetes:///a.b.svc.cluster.local:port", targetInfo{"a", "b", "port", true, false}, false},
	} {
		got, err := parseResolverTarget(parseTarget(test.target))
		if err == nil && test.err {
			t.Errorf("case %d: want error but got nil", i)
			continue
		}
		if err != nil && !test.err {
			t.Errorf("case %d:got '%v' error but don't want an error", i, err)
			continue
		}
		if got != test.want {
			t.Errorf("case %d: parseTarget(%q) = %+v, want %+v", i, test.target, got, test.want)
		}
	}
}

func parseTarget(target string) resolver.Target {
	u, err := url.Parse(target)
	if err != nil {
		return resolver.Target{}
	}
	endpoint := u.Path
	if endpoint == "" {
		endpoint = u.Opaque
	}
	endpoint = strings.TrimPrefix(endpoint, "/")

	// return resolver.Target{
	// 	Scheme:    u.Scheme,
	// 	Authority: u.Host,
	// 	Endpoint:  endpoint,
	// }
	return resolver.Target{
		URL: url.URL{
			Scheme: u.Scheme,
			Host:   u.Host,
			Path:   endpoint,
		},
	}
}
