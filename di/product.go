package di

import (
	"github.com/gabriwl165/clean-arch-go/adapter/http/productservice"
	"github.com/gabriwl165/clean-arch-go/adapter/postgres"
	"github.com/gabriwl165/clean-arch-go/adapter/postgres/productrepository"
	"github.com/gabriwl165/clean-arch-go/core/domain"
	"github.com/gabriwl165/clean-arch-go/core/usecase/productusecase"
)

func ConfigProductDI(conn postgres.PoolInterface) domain.ProductService {
	productRepository := productrepository.New(conn)
	productUseCase := productusecase.New(productRepository)
	ProductService := productservice.New(productUseCase)
	return ProductService
}
