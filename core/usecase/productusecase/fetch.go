package productusecase

import (
	"github.com/gabriwl165/clean-arch-go/core/domain"
	"github.com/gabriwl165/clean-arch-go/core/dto"
)

func (usecase usecase) Fetch(paginationRequest *dto.PaginationRequestParams) (*domain.Pagination[[]domain.Product], error) {
	products, err := usecase.repository.Fetch(paginationRequest)
	if err != nil {
		return nil, err
	}
	return products, nil
}
