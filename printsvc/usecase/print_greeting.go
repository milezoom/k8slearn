package usecase

import (
	"context"
	"fmt"
	hContract "hellocln/contract"
	"printsvc/contract"
)

func (uc *UseCase) PrintGreeting(ctx context.Context, request *contract.EmptyRequest) (response *contract.PrintGreetingResponse, err error) {
	resp, err := uc.Repo.HelloSvc.PrintHello(ctx, &hContract.PrintHelloRequest{
		Name: "World",
	})
	if err != nil {
		return nil, err
	}
	return &contract.PrintGreetingResponse{
		Message: fmt.Sprintf("%s\nGood Day Today~", resp.GetMessage()),
	}, nil
}
