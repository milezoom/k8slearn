// Code generated by proto-gen-go-lib-contract . DO NOT EDIT.
// Source: "hellocln/contract"

package implementor

import (
	"context"
	"hellocln/contract"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type HelloServiceImplementor struct {
	GrpcConn *grpc.ClientConn
	Cli      contract.HelloServiceClient
}

func (i *HelloServiceImplementor) Close() error {
	if i.GrpcConn != nil {
		return i.GrpcConn.Close()
	}
	return nil

}

// TODO: update function below
func (i *HelloServiceImplementor) GetDefaultMetadata(ctx context.Context) metadata.MD {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return md
	} else {
		return metadata.MD{}
	}
}
