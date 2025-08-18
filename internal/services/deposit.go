package services

import (
	"context"

	"github.com/wafi04/otomaxv2/internal/integrations/duitku"
	"github.com/wafi04/otomaxv2/internal/model"
	"github.com/wafi04/otomaxv2/internal/repository"
)

type DepositService struct {
	repo   *repository.DepositRepository
	duitku *duitku.DuitkuService
}

func NewDuitkuService(repo *repository.DepositRepository, duitku *duitku.DuitkuService) *DepositService {
	return &DepositService{
		repo:   repo,
		duitku: duitku,
	}
}

func (ds *DepositService) Create(c context.Context, req model.CreateDeposit) (bool, error) {
	return ds.repo.Create(c, req)
}

func (ds *DepositService) GetAll(c context.Context, req model.FilterDeposit) ([]model.DepositData, int, error) {
	return ds.repo.GetAll(c, req)
}

func (ds *DepositService) GetByID(c context.Context, id int) (*model.DepositData, error) {
	return ds.repo.GetByID(c, id)
}
