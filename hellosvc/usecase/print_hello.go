package usecase

import (
	"context"
	"fmt"
	"hellosvc/contract"
)

func (*UseCase) PrintHello(ctx context.Context, request *contract.PrintHelloRequest) (response *contract.PrintHelloResponse, err error) {
	return &contract.PrintHelloResponse{
		Message: fmt.Sprintf("Hello, %s!", request.GetName()),
	}, nil
}
