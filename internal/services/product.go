package services

import (
	"github.com/gin-gonic/gin"
	"github.com/wafi04/otomaxv2/internal/model"
	"github.com/wafi04/otomaxv2/internal/repository"
)


type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}



func (r *ProductService) GetAll(c *gin.Context)([]*model.ProductData, error){
	return r.productRepo.GetProducts(c,10,1)
}