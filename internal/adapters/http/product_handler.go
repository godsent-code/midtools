package http

import (
	"net/http"

	"github.com/godsent-code/midtools/internal/application/product"
	"github.com/godsent-code/midtools/pkg"
)

type ProductHandler struct {
	service product.ProductService
}

func (ph *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := ph.service.GetProducts(r.Context())
	if err != nil {
		pkg.WriteResponse(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	pkg.WriteResponse(w, http.StatusOK, products)
}

func (ph *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	err := ph.service.CreateProducts(r.Context())
	if err != nil {
		pkg.WriteResponse(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	pkg.WriteResponse(w, http.StatusOK, "Product created")
}

func NewProductHandler(service product.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}
