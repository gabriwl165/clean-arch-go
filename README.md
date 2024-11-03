# A Simple Implementation Clean Architecture in Go

In this guide, i'm going to explain how the project works

We currently have 4 folders:
- adapter/
- core/
- database/
- di/

## Adapter

### In this folder, we've placed our adapters for external services, such as HTTP and PostgreSQL.
Let's imagine we need to change how we interact with the PostgreSQL database, or perhaps we want to switch the framework we're using to handle requests. By creating a boundary between our software and external dependencies, we make these changes easier to manage.

- postgres/

In this folder we've created our `connector.go` (`adapter/postgres/connector.go`), this file is simple, basically he's implementing our connection into the database and providing some function that will interact with him

At the `new.go` file (`adapter/postgres/productrepository/new.go`), we're creating a struct like this and a function called `New`:
```go
type repository struct {
	db postgres.PoolInterface
}

func New(db postgres.PoolInterface) domain.ProductRepository {
	return &repository{
		db: db,
	}
}
```
The New function indicates that it will return a `domain.ProductRepository`, which is an interface.
```go
type ProductRepository interface {
	Create(productRequest *dto.CreateProductRequest) (*Product, error)
	Fetch(paginationRequest *dto.PaginationRequestParams) (*Pagination[[]Product], error)
}
```
This interface indicates that it will return an object with two methods, `Create` and `Fetch`. 
However, our repository struct does not implement them—not in this file, at least. 
Since `create`, `fetch`, and `new` are in the same `package`, those other methods are implemented in separate files. 

In our fetch.go file, we’ve begun instantiating some variables:
```go
ctx := context.Background()
products := []domain.Product{}
total := int32(0)
```
- `ctx` is receiving the context, which will be used shortly.
- `products` is our array that stores all the products fetched from the database and will be returned to the client.
- `total`, represents the total number of records, since our API is paginated. We need to inform the client how much data is available to be fetched.

```go
	query, queryCount, err := paginate.Paginate("SELECT * FROM product").
		Page(pagination.Page).
		Desc(pagination.Descending).
		Sort(pagination.Sort).
		RowsPerPage(pagination.ItemsPerPage).
		SearchBy(pagination.Search, "name", "description").
		Query()
```

By using `paginate.Paginate`, we're fetching data from the database, applying filters based on the data provided in the request, since our method is receiving the pagination parameters:

`func (repository repository) Fetch(pagination *dto.PaginationRequestParams) (*domain.Pagination[[]domain.Product], error)`

After executing the query, we're simply converting the data and determining how many rows exist in the database that match the request parameters.
```go
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
```
In our `create.go` file, we’re simply inserting the data into the database and validating whether everything worked correctly. 
There's nothing to worry about.
```go
func (repository repository) Create(productRequest *dto.CreateProductRequest) (*domain.Product, error) {
	ctx := context.Background()
	product := domain.Product{}
	err := repository.db.QueryRow(
		ctx,
		"INSERT INTO product (name, price, description) VALUES ($1, $2, $3) returning *",
		productRequest.Name,
		productRequest.Price,
		productRequest.Description,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Description,
	)

	if err != nil {
		return nil, err
	}

	return &product, nil

}
```

With all the necessary implementations for the `domain.ProductRepository`, we are now able to respond with the following:
```go
func New(db postgres.PoolInterface) domain.ProductRepository {
	return &repository{
		db: db,
	}
}
```

- http/

In this folder, we will implement the logic for handling requests from our endpoint. 
In the `New` function (`adapter/http/productservice/new.go`), we have a service struct and the New function:

```go
type service struct {
	usecase domain.ProductUseCase
}

func New(usecase domain.ProductUseCase) domain.ProductService {
	return &service{
		usecase: usecase,
	}
}
```
Our function returns an instance of the `service struct`, but we're expecting a `domain.ProductService`, which requires that this method be implemented in our object:
```go
type ProductService interface {
	Create(response http.ResponseWriter, request *http.Request)
	Fetch(response http.ResponseWriter, request *http.Request)
}
```
This is implemented in our `create.go` (`adapter/postgres/productrepository/create.go`) and `fetch.go` (`adapter/postgres/productrepository/fetch.go`) files. 
Both files have essentially the same implementation, with the difference being what is delegated to the database.

In `create.go`, we use the use cases to persist (`Create`) the data into our database.
```go
func (service service) Create(response http.ResponseWriter, request *http.Request) {
	productRequest, err := dto.FromJSONCreateProductRequest(request.Body)

	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(err.Error()))
	}

	product, err := service.usecase.Create(productRequest)
	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(err.Error()))
	}

	json.NewEncoder(response).Encode(product)
}
```

The request comes in, is parsed into a `productRequest`, and then delegated to the use cases of the service, which are specified in our service struct:
```
type service struct {
	usecase domain.ProductUseCase
}
```

The difference in `fetch.go` is that we call `products, err := service.usecase.Fetch(paginationRequest)`. 
With everything implemented, our New function will work correctly.
```go
func New(db postgres.PoolInterface) domain.ProductRepository {
	return &repository{
		db: db,
	}
}
```
## Core

### In this folder, we’ve placed what’s referred to as the domain of our application, which includes business logic like use cases and DTOs.

- domain/

In our domain, we've implemented a generic `Pagination struct`:
```go
type Pagination[T any] struct {
	Items T     `json:"items"`
	Total int32 `json:"total"`
}
```
And in our `product.go` file, we have the full definition of the interface that we've been using in our `adapter` folder.
```go
type Product struct {
	ID          int32   `json:"id"`
	Name        string  `json:"name"`
	Price       float32 `json:"price"`
	Description string  `json:"description"`
}

type ProductService interface {
	Create(response http.ResponseWriter, request *http.Request)
	Fetch(response http.ResponseWriter, request *http.Request)
}

type ProductUseCase interface {
	Create(productRequest *dto.CreateProductRequest) (*Product, error)
	Fetch(paginationRequest *dto.PaginationRequestParams) (*Pagination[[]Product], error)
}

type ProductRepository interface {
	Create(productRequest *dto.CreateProductRequest) (*Product, error)
	Fetch(paginationRequest *dto.PaginationRequestParams) (*Pagination[[]Product], error)
}

```

- dto/

In this folder, we have the JSON that will be received in the request and the JSON that will be sent as the response to the client. 
It's as simple as that.

- usecase/

In the `usecases` folder, we will place all our business logic, such as validation, conversions, calls to our adapters, and more. 
In our `new.go` file (`core/usecase/productusecase/new.go`):
```go
type usecase struct {
	repository domain.ProductRepository
}

func New(repository domain.ProductRepository) domain.ProductUseCase {
	return &usecase{
		repository: repository,
	}
}
```
Our `New` function returns a `usecase struct`, but the return type must be `domain.ProductUseCase`, which is defined in our `core/domain`.
```go
type ProductUseCase interface {
	Create(productRequest *dto.CreateProductRequest) (*Product, error)
	Fetch(paginationRequest *dto.PaginationRequestParams) (*Pagination[[]Product], error)
}
```
So, in our `create` and `fetch` files, we’re implementing those methods, which call our `domain.ProductRepository` to perform CRUD operations on the database.

## Database

### In this folder, we’ve placed our migrations, which will be executed when postgres.RunMigrations() is called.

## DI

### In this folder, we’ve placed our factory, which creates the ProductService that will be used by our endpoint.

Here, we have our factory that configures the dependencies needed for our route to function.
```go
func ConfigProductDI(conn postgres.PoolInterface) domain.ProductService {
	productRepository := productrepository.New(conn)
	productUseCase := productusecase.New(productRepository)
	ProductService := productservice.New(productUseCase)
	return ProductService
}
```

Just as a reminder, ConfigProductDI is being called in our `main` file (`adapter/http/main.go`)
```go
func main() {
	ctx := context.Background()
	conn := postgres.GetConnection(ctx)
	defer conn.Close()

	postgres.RunMigrations()
	.....
```