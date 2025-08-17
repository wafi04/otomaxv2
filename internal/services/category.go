package services

import (
	"context"

	"github.com/wafi04/otomaxv2/internal/model"
	"github.com/wafi04/otomaxv2/internal/repository"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) CreateCategory(ctx context.Context, input model.CreateCategory) error {
	return s.categoryRepo.Create(ctx, input)
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id int) (*model.Category, error) {
	return s.categoryRepo.GetByID(ctx, id)
}

func (s *CategoryService) GetCategoryByCode(ctx context.Context, code string, subCategoryId *int) (*model.CategoryCodeResponse, error) {
	return s.categoryRepo.GetByCodeWithFilter(ctx, code, &repository.CategoryFilter{
		SubCategoryID: subCategoryId,
	})
}

func (s *CategoryService) GetAllCategories(ctx context.Context, skip, limit int, search, filterType string, active string) ([]model.Category, int, error) {
	return s.categoryRepo.GetAll(ctx, skip, limit, search, filterType, active)
}
func (s *CategoryService) UpdateCategory(ctx context.Context, id int, input model.CreateCategory) error {
	return s.categoryRepo.Update(ctx, id, input)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	return s.categoryRepo.Delete(ctx, id)
}

func (repo *CategoryService) Count(ctx context.Context, search, filterType string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM categories
		WHERE ($1 = '' OR name ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR type = $2)
	`
	var total int
	err := repo.categoryRepo.DB.QueryRowContext(ctx, query, search, filterType).Scan(&total)
	return total, err
}
