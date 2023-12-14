package implementor

import (
	"context"

	"hellocln/contract"
	"google.golang.org/grpc/metadata"
)

func (i *HelloServiceImplementor) PrintHello(ctx context.Context, md metadata.MD, input *contract.PrintHelloRequest) (resp *contract.PrintHelloResponse, err error) {

	// build metadata
	ctx = metadata.NewOutgoingContext(ctx, md)

	return i.Cli.PrintHello(ctx, input)
}
