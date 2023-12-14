package usecase

import (
	"printsvc/repository"
)

var useCasePointer *UseCase

type UseCase struct {
	Repo *repository.Repository
}

func NewUsecase(repoIn *repository.Repository) *UseCase {
	if useCasePointer == nil {
		useCasePointer = &UseCase{Repo: repoIn}
	}
	return useCasePointer
}
