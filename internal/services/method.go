package services

import (
	"context"

	"github.com/wafi04/otomaxv2/internal/model"
	"github.com/wafi04/otomaxv2/internal/repository"
)

type MethodService struct {
	Repo *repository.MethodRepository
}

func NewMethodService(Repo *repository.MethodRepository) *MethodService {
	return &MethodService{
		Repo: Repo,
	}
}

func (service *MethodService) Create(c context.Context, data model.CreateMethodData) (*model.MethodData, error) {
	return service.Repo.Create(c, &data)
}

func (service *MethodService) GetAll(c context.Context, skip, limit int, search, filterType string, active string) ([]model.MethodData, int, error) {
	return service.Repo.GetAll(c, skip, limit, search, filterType, active)

}

func (service *MethodService) GetAllGroupedByType(c context.Context) ([]repository.MethodGroupResponse, error) {
	return service.Repo.GetAllGroupedByType(c)

}

func (service *MethodService) Update(c context.Context, id int, data model.UpdateMethodData) (*model.MethodData, error) {
	return service.Repo.Update(c, id, &data)
}

func (service *MethodService) Delete(c context.Context, id int) error {
	return service.Repo.Delete(c, id)
}
