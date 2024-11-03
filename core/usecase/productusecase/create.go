package productusecase

import (
	"github.com/gabriwl165/clean-arch-go/core/domain"
	"github.com/gabriwl165/clean-arch-go/core/dto"
)

func (usecase usecase) Create(productRequest *dto.CreateProductRequest) (*domain.Product, error) {
	product, err := usecase.repository.Create(productRequest)
	if err != nil {
		return nil, err
	}

	return product, err
}
