package hellosvc

import (
	"context"
	"hellocln/contract"
)

type HelloSvcClient struct {
	cli contract.HelloService
}

func (c *HelloSvcClient) PrintHello(ctx context.Context, input *contract.PrintHelloRequest) (*contract.PrintHelloResponse, error) {
	md := c.cli.GetDefaultMetadata(ctx)
	return c.cli.PrintHello(ctx, md, input)
}
