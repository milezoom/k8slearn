# iresolver

GRPC Kubernetes Resolver

![ALT](/diagram.png)


## How to Use

### Just add this before GRPC dial    
```
iresolver.RegisterInCluster() 
```

### Address should like this
```
kubernetes:///supply-prediction:9102
kubernetes:///supply-prediction.default:9102
```

### Add Service Account in Deployment.yaml
```
apiVersion: apps/v1 
kind: Deployment 
metadata: 
  labels: 
    app: heatmap-prediction 
  name: heatmap-prediction 
spec: 
  replicas: 7 
  selector: 
    matchLabels: 
      app: heatmap-prediction 
  strategy: 
    rollingUpdate: 
      # This should be set to 0 if we only have 1 replica defined 
      maxUnavailable: 25% 
  template: 
    metadata: 
      labels: 
        app: heatmap-prediction 
    spec: 
      serviceAccountName: heatmap-prediction-sa # service account
```

### Using Headless Service

```
apiVersion: v1
kind: Service
metadata:
  name: heatmap-prediction
  labels:
    app: heatmap-prediction
  annotations:
    cloud.google.com/app-protocols: '{"grpc":"HTTP2"}'
    beta.cloud.google.com/backend-config: '{"ports": {"9100":"heatmap-prediction-config"}}'
spec:
  type: NodePort
  ports:
    - name: grpc
      port: 9100
      targetPort: 9100
  selector:
    app: heatmap-prediction
```

### Example Create Connection
```
func grpcConnectionV2(address string, creds credentials.TransportCredentials) (*grpc.ClientConn, error) { 
	var conn *grpc.ClientConn 
	var err error 
	var opts []grpc.DialOption 
	opts = append( 
		opts, 
		grpc.WithUnaryInterceptor(apmgrpc.NewUnaryClientInterceptor()), 
		// grpc.WithReturnConnectionError(), 
		grpc.WithKeepaliveParams(keepalive.ClientParameters{ 
			// Time:                time.Duration(*flagKeepAliveTime) * time.Second, 
			Timeout:             time.Duration(1) * time.Second, 
			PermitWithoutStream: true, 
		}), 
	) 
	if strings.Contains(address, "dns:///") { 
		resolver.Register(dns.NewBuilder()) 
		opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`)) 
	} 
	if strings.Contains(address, "kubernetes:///") { 
		iresolver.RegisterInCluster() 
		opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`)) 
	} 
	if creds == nil { 
		opts = append(opts, grpc.WithInsecure()) 
		conn, err = grpc.Dial(address, opts...) 
	} else { 
		opts = append(opts, grpc.WithTransportCredentials(creds)) 
		conn, err = grpc.Dial(address, opts...) 
	} 
	if err != nil { 
		return nil, err 
	} 
	return conn, nil 
} 
```