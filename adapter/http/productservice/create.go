package productservice

import (
	"encoding/json"
	"net/http"

	"github.com/gabriwl165/clean-arch-go/core/dto"
)

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
