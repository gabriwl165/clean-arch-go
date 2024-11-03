package productrepository

import (
	"github.com/gabriwl165/clean-arch-go/adapter/postgres"
	"github.com/gabriwl165/clean-arch-go/core/domain"
)

type repository struct {
	db postgres.PoolInterface
}

func New(db postgres.PoolInterface) domain.ProductRepository {
	return &repository{
		db: db,
	}
}
