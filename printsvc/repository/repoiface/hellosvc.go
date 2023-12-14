package repoiface

import (
	"context"
	"hellocln/contract"
)

type HelloSvc interface {
	PrintHello(ctx context.Context, input *contract.PrintHelloRequest) (*contract.PrintHelloResponse, error)
}
