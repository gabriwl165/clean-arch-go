package productusecase

import "github.com/gabriwl165/clean-arch-go/core/domain"

type usecase struct {
	repository domain.ProductRepository
}

func New(repository domain.ProductRepository) domain.ProductUseCase {
	return &usecase{
		repository: repository,
	}
}
