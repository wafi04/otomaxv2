package services

import (
	"github.com/wafi04/otomaxv2/internal/model"
	"github.com/wafi04/otomaxv2/internal/repository"
)

type NewsService struct {
	newsRepo *repository.NewsRepository
}

func NewNewsService(newsRepo *repository.NewsRepository) *NewsService {
	return &NewsService{
		newsRepo: newsRepo,
	}
}

func (service *NewsService) Create(req model.CreateNews) (*model.News, error) {
	return service.newsRepo.Create(&req)
}

func (service *NewsService) GetAll(status, newsType *string) ([]model.News, error) {
	return service.newsRepo.GetAll(status, newsType)
}

func (service *NewsService) GetByID(id int) (*model.News, error) {
	return service.newsRepo.GetByID(id)
}

func (service *NewsService) Update(id int, req model.CreateNews) (*model.News, error) {
	return service.newsRepo.Update(id, &req)
}

func (service *NewsService) Delete(id int) error {
	return service.newsRepo.Delete(id)
}