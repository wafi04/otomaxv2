package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/wafi04/otomaxv2/internal/model"
)

type CategoryRepository struct {
	DB *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

// Create
func (repo *CategoryRepository) Create(ctx context.Context, category model.CreateCategory) error {
	query := `
		INSERT INTO categories (
			name, sub_name, brand, code, is_check_nickname, status,
			thumbnail, type, instruction, information, banner, placeholder_1, placeholder_2,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13,
			NOW(), NOW()
		)
	`

	_, err := repo.DB.ExecContext(ctx, query,
		category.Name, category.SubName, category.Brand, category.Code,
		category.IsCheckNickname, category.Status, category.Thumbnail, category.Type,
		category.Instruction, category.Information, category.Banner,
		category.Placeholder1, category.Placeholder2,
	)
	if err != nil {
		log.Printf("Create Category error: %v", err)
	}
	return err
}

func (repo *CategoryRepository) Update(ctx context.Context, id int, category model.CreateCategory) error {
	query := `
		UPDATE categories
		SET name = $1, sub_name = $2, brand = $3, code = $4, is_check_nickname = $5, status = $6,
			thumbnail = $7, type = $8, instruction = $9, information = $10,
			banner = $11, placeholder_1 = $12, placeholder_2 = $13,
			updated_at = NOW()
		WHERE id = $14
	`

	_, err := repo.DB.ExecContext(ctx, query,
		category.Name, category.SubName, category.Brand, category.Code,
		category.IsCheckNickname, category.Status, category.Thumbnail, category.Type,
		category.Instruction, category.Information, category.Banner,
		category.Placeholder1, category.Placeholder2,
		id,
	)
	if err != nil {
		log.Printf("Update Category error: %v", err)
	}
	return err
}

// Delete
func (repo *CategoryRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1`

	_, err := repo.DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("Delete Category error: %v", err)
	}
	return err
}
func (repo *CategoryRepository) GetAll(ctx context.Context, skip, limit int, search, filterType string, active string) ([]model.Category, int, error) {
	// Query untuk mendapatkan total count
	countQuery := `
		SELECT COUNT(*) 
		FROM categories
		WHERE ($1 = '' OR name ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR type = $2)
		  AND ($3 = '' OR status = $3)
	`

	var totalCount int
	err := repo.DB.QueryRowContext(ctx, countQuery, search, filterType, active).Scan(&totalCount)
	if err != nil {
		log.Printf("GetAll Categories count error: %v", err)
		return nil, 0, err
	}

	query := `
		SELECT id, name, sub_name, brand, code, is_check_nickname, status,
			thumbnail, type, instruction, information,
			banner, placeholder_1, placeholder_2, created_at, updated_at
		FROM categories
		WHERE ($1 = '' OR name ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR type = $2)
		  AND ($3 = '' OR status = $3)
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5
	`

	rows, err := repo.DB.QueryContext(ctx, query, search, filterType, active, limit, skip)
	if err != nil {
		log.Printf("GetAll Categories error: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var cat model.Category
		err := rows.Scan(
			&cat.ID, &cat.Name, &cat.SubName, &cat.Brand, &cat.Code,
			&cat.IsCheckNickname, &cat.Status, &cat.Thumbnail, &cat.Type,
			&cat.Instruction, &cat.Information, &cat.Banner, &cat.Placeholder1,
			&cat.Placeholder2, &cat.CreatedAt, &cat.UpdatedAt,
		)
		if err != nil {
			log.Printf("Scan Category error: %v", err)
			continue
		}
		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return categories, totalCount, nil
}

func (repo *CategoryRepository) GetByID(ctx context.Context, id int) (*model.Category, error) {
	query := `
		SELECT id, name, sub_name, brand, code, is_check_nickname, status,
			thumbnail, type, instruction, information,
			banner, placeholder_1, placeholder_2, created_at, updated_at
		FROM categories
		WHERE id = $1
	`

	var cat model.Category
	err := repo.DB.QueryRowContext(ctx, query, id).Scan(
		&cat.ID, &cat.Name, &cat.SubName, &cat.Brand, &cat.Code,
		&cat.IsCheckNickname, &cat.Status, &cat.Thumbnail, &cat.Type,
		&cat.Instruction, &cat.Information, &cat.Banner, &cat.Placeholder1,
		&cat.Placeholder2, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("GetByID Category error: %v", err)
		return nil, err
	}
	return &cat, nil
}

type CategoryFilter struct {
	SubCategoryID *int
}

func (repo *CategoryRepository) GetByCodeWithFilter(ctx context.Context, code string, filter *CategoryFilter) (*model.CategoryCodeResponse, error) {
	// Query utama untuk category
	categoryQuery := `
		SELECT 
			id, name, sub_name, brand, code, is_check_nickname, 
			status, thumbnail, type, instruction, information, banner
		FROM categories 
		WHERE code = $1`

	var cat model.CategoryCodeResponse
	err := repo.DB.QueryRowContext(ctx, categoryQuery, code).Scan(
		&cat.ID, &cat.Name, &cat.SubName, &cat.Brand, &cat.Code,
		&cat.IsCheckNickname, &cat.Status, &cat.Thumbnail, &cat.Type,
		&cat.Instruction, &cat.Information, &cat.Banner,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // category tidak ditemukan
		}
		return nil, err
	}

	// Query untuk subcategories yang terkait dengan category ini
	subCategoryQuery := `
		SELECT id, name, code, is_active
		FROM sub_categories 
		WHERE category_id = $1 AND is_active = 'active'
		ORDER BY name`

	subRows, err := repo.DB.QueryContext(ctx, subCategoryQuery, cat.ID)
	if err != nil {
		fmt.Print("err", err)
		return nil, err
	}
	defer subRows.Close()

	for subRows.Next() {
		var subCat model.SubCategory
		err := subRows.Scan(&subCat.Id, &subCat.Name, &subCat.Code, &subCat.Status)
		if err != nil {
			return nil, err
		}
		cat.SubCategories = append(cat.SubCategories, subCat)
	}

	if err = subRows.Err(); err != nil {
		return nil, err
	}

	var productQuery string
	var productArgs []interface{}

	if filter != nil && filter.SubCategoryID != nil && *filter.SubCategoryID > 0 {
		productQuery = `
			SELECT id, name, price, denomination_type, sub_category_id
			FROM products 
			WHERE category_id = $1 AND sub_category_id = $2 AND status = 'active'
			ORDER BY name`
		productArgs = []interface{}{cat.ID, *filter.SubCategoryID}
	} else {
		// Jika tidak ada filter subcategory, ambil semua products dari category
		productQuery = `
			SELECT id, name, price, denomination_type, sub_category_id
			FROM products 
			WHERE category_id = $1 AND status = 'active'
			ORDER BY name`
		productArgs = []interface{}{cat.ID}
	}

	productRows, err := repo.DB.QueryContext(ctx, productQuery, productArgs...)
	if err != nil {
		return nil, err
	}
	defer productRows.Close()

	for productRows.Next() {
		var prod model.Product
		err := productRows.Scan(
			&prod.ID, &prod.Name, &prod.Price,
			&prod.DenominationType, &prod.SubCategoryID,
		)
		if err != nil {
			return nil, err
		}
		cat.Products = append(cat.Products, prod)
	}

	if err = productRows.Err(); err != nil {
		return nil, err
	}

	return &cat, nil
}
func (repo *CategoryRepository) Count(ctx context.Context, search, filterType string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM categories
		WHERE ($1 = '' OR name ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR type = $2)
	`
	var total int
	err := repo.DB.QueryRowContext(ctx, query, search, filterType).Scan(&total)
	if err != nil {
		log.Printf("Count Category error: %v", err)
		return 0, err
	}
	return total, nil
}
