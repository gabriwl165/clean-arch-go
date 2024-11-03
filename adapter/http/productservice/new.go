package productservice

import "github.com/gabriwl165/clean-arch-go/core/domain"

type service struct {
	usecase domain.ProductUseCase
}

func New(usecase domain.ProductUseCase) domain.ProductService {
	return &service{
		usecase: usecase,
	}
}
