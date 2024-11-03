package productrepository

import (
	"context"

	"github.com/booscaaa/go-paginate/paginate"
	"github.com/gabriwl165/clean-arch-go/core/domain"
	"github.com/gabriwl165/clean-arch-go/core/dto"
)

func (repository repository) Fetch(pagination *dto.PaginationRequestParams) (*domain.Pagination[[]domain.Product], error) {
	ctx := context.Background()
	products := []domain.Product{}
	total := int32(0)

	query, queryCount, err := paginate.Paginate("SELECT * FROM product").
		Page(pagination.Page).
		Desc(pagination.Descending).
		Sort(pagination.Sort).
		RowsPerPage(pagination.ItemsPerPage).
		SearchBy(pagination.Search, "name", "description").
		Query()

	if err != nil {
		return nil, err
	}
	{
		rows, err := repository.db.Query(
			ctx, *query,
		)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			product := domain.Product{}
			rows.Scan(
				&product.ID,
				&product.Name,
				&product.Price,
				&product.Description,
			)
			products = append(products, product)
		}
	}
	{
		err := repository.db.QueryRow(ctx, *queryCount).Scan(&total)
		if err != nil {
			return nil, err
		}
	}
	return &domain.Pagination[[]domain.Product]{
		Items: products,
		Total: total,
	}, nil

}
