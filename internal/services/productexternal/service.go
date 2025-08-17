package productexternal

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/wafi04/otomaxv2/internal/integrations/digiflazz"
)

type ProductExternal struct {
	DigiflazzService *digiflazz.DigiflazzService
	DB               *sql.DB
}

func NewProductExternal(digiService *digiflazz.DigiflazzService, db *sql.DB) *ProductExternal {
	return &ProductExternal{
		DigiflazzService: digiService,
		DB:               db,
	}
}

func (repo *ProductExternal) getCategoryAndSubCategory(ctx context.Context, categoryName, brand string) (int, int, error) {
	var categoryID sql.NullInt64

	queryCat := `SELECT id FROM categories WHERE LOWER(brand) = LOWER($1)`
	err := repo.DB.QueryRowContext(ctx, queryCat, brand).Scan(&categoryID)
	if err != nil {
		queryCatName := `SELECT id FROM categories WHERE LOWER(name) = LOWER($1)`
		err = repo.DB.QueryRowContext(ctx, queryCatName, categoryName).Scan(&categoryID)
		if err != nil {
			return 0, 0, fmt.Errorf("category not found for brand: %s, category: %s", brand, categoryName)
		}
	}
	return int(categoryID.Int64), 1, nil
}

// FIXED: Method untuk process products dari Digiflazz
func (pe *ProductExternal) GetProductDigiflazz(ctx context.Context, data []*digiflazz.InternalProduct) ([]*digiflazz.InternalProduct, error) {
	var processedProducts []*digiflazz.InternalProduct

	for _, product := range data {
		// Get category and subcategory for each product
		categoryID, subCategoryID, err := pe.getCategoryAndSubCategory(ctx, product.Category, product.Brand)
		if err != nil {
			log.Printf("Warning: Failed to get category for product %s: %v", product.ProviderCode, err)
			continue
		}

		processedProduct := &digiflazz.InternalProduct{
			ProviderCode:  product.ProviderCode,
			ProviderName:  product.ProviderName,
			Category:      product.Category,
			Brand:         product.Brand,
			Type:          product.Type,
			Description:   product.Description,
			CostPrice:     product.CostPrice,
			SellingPrice:  product.SellingPrice,
			ProfitMargin:  product.ProfitMargin,
			Stock:         product.Stock,
			IsUnlimited:   product.IsUnlimited,
			IsActive:      product.IsActive,
			Status:        product.Status,
			SellerName:    product.SellerName,
			StartCutOff:   product.StartCutOff,
			EndCutOff:     product.EndCutOff,
			SupportMulti:  product.SupportMulti,
			Provider:      product.Provider,
			CategoryID:    &categoryID,
			SubCategoryID: &subCategoryID,
		}

		// Optional: Save to database
		err = pe.SaveProducts(ctx, processedProduct)
		if err != nil {
			log.Printf("Warning: Failed to save product %s: %v", product.ProviderCode, err)
			// Continue processing other products
		}

		processedProducts = append(processedProducts, processedProduct)
	}

	return processedProducts, nil
}

func (pe *ProductExternal) SaveProducts(ctx context.Context, product *digiflazz.InternalProduct) error {
	log.Printf("=== PROCESSING PRODUCT: %s ===", product.ProviderCode)

	// Debug overflow check (dari kode sebelumnya)
	const maxInt32 = 2147483647
	const minInt32 = -2147483648

	if product.CostPrice > maxInt32 || product.CostPrice < minInt32 {
		log.Printf("OVERFLOW: CostPrice %s", product.ProviderCode)
		return fmt.Errorf("cost price overflow for product %s", product.ProviderCode)
	}
	if product.SellingPrice > maxInt32 || product.SellingPrice < minInt32 {
		log.Printf("OVERFLOW: SellingPrice %s", product.ProviderCode)
		return fmt.Errorf("selling price overflow for product %s", product.ProviderCode)
	}

	// Mulai transaction untuk konsistensi data
	tx, err := pe.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// 1. Save/Update Provider Product
	err = pe.saveOrUpdateProviderProduct(ctx, tx, product)
	if err != nil {
		return fmt.Errorf("failed to save provider product: %v", err)
	}

	// 2. Save/Update Main Product
	err = pe.saveOrUpdateMainProduct(ctx, tx, product)
	if err != nil {
		return fmt.Errorf("failed to save main product: %v", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("=== SUCCESS PROCESSING: %s ===", product.ProviderCode)
	return nil
}

// Helper function untuk save/update provider_products
func (pe *ProductExternal) saveOrUpdateProviderProduct(ctx context.Context, tx *sql.Tx, product *digiflazz.InternalProduct) error {
	// Check if provider product exists
	var existingID int
	checkQuery := `
		SELECT id FROM provider_products 
		WHERE provider_code = $1 AND provider_id = (SELECT id FROM providers WHERE slug = $2)`

	err := tx.QueryRowContext(ctx, checkQuery, product.ProviderCode, product.Provider).Scan(&existingID)

	if err == sql.ErrNoRows {
		// Insert new provider product
		log.Printf("Inserting new provider product: %s", product.ProviderCode)

		insertQuery := `
			INSERT INTO provider_products 
			(provider_id, provider_code, provider_name, cost_price, selling_price, 
			 profit_margin, stock, status, is_available, created_at, updated_at)
			VALUES 
			((SELECT id FROM providers WHERE slug = $1), $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())`

		_, err = tx.ExecContext(ctx, insertQuery,
			product.Provider,
			product.ProviderCode,
			product.ProviderName,
			product.CostPrice,
			product.SellingPrice,
			product.ProfitMargin,
			product.Stock,
			product.Status,
			product.IsActive)

		if err != nil {
			log.Printf("INSERT ERROR provider product %s: %v", product.ProviderCode, err)
			return err
		}
	} else if err != nil {
		return fmt.Errorf("failed to check existing provider product: %v", err)
	} else {
		// Update existing provider product
		log.Printf("Updating existing provider product: %s (ID: %d)", product.ProviderCode, existingID)

		updateQuery := `
			UPDATE provider_products 
			SET provider_name = $1, cost_price = $2, selling_price = $3, 
				profit_margin = $4, stock = $5, status = $6, is_available = $7, 
				updated_at = NOW()
			WHERE id = $8`

		_, err = tx.ExecContext(ctx, updateQuery,
			product.ProviderName,
			product.CostPrice,
			product.SellingPrice,
			product.ProfitMargin,
			product.Stock,
			product.Status,
			product.IsActive,
			existingID)

		if err != nil {
			log.Printf("UPDATE ERROR provider product %s: %v", product.ProviderCode, err)
			return err
		}
	}

	return nil
}

func (pe *ProductExternal) saveOrUpdateMainProduct(ctx context.Context, tx *sql.Tx, product *digiflazz.InternalProduct) error {
	// Check if main product exists (berdasarkan provider_code atau kriteria lain)
	var existingProductID int
	checkMainQuery := `
		SELECT id FROM products 
		WHERE name = $1 OR description LIKE $2`

	searchPattern := "%" + product.ProviderCode + "%"
	err := tx.QueryRowContext(ctx, checkMainQuery, product.ProviderName, searchPattern).Scan(&existingProductID)

	if err == sql.ErrNoRows {
		// Insert new main product
		log.Printf("Inserting new main product for: %s", product.ProviderCode)

		// Mapping category berdasarkan provider atau product type

		insertMainQuery := `
			INSERT INTO products (
				category_id, sub_category_id, name, description, price, original_price,
				denomination, denomination_type, sort_order, status, stock, 
				created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())`

		_, err = tx.ExecContext(ctx, insertMainQuery,
			product.CategoryID,
			product.SubCategoryID,
			product.ProviderName,
			fmt.Sprintf("Provider: %s, Code: %s", product.Provider, product.ProviderCode),
			product.SellingPrice,
			product.CostPrice,
			pe.getDenomination(product),     // Extract dari product name/code
			pe.getDenominationType(product), // Mobile, Game, etc
			pe.getSortOrder(product),        // Berdasarkan kategori
			product.Status,
			product.Stock)

		if err != nil {
			log.Printf("INSERT ERROR main product %s: %v", product.ProviderCode, err)
			return err
		}
	} else if err != nil {
		return fmt.Errorf("failed to check existing main product: %v", err)
	} else {
		// Update existing main product
		log.Printf("Updating existing main product: %s (ID: %d)", product.ProviderName, existingProductID)

		updateMainQuery := `
			UPDATE products 
			SET price = $1, original_price = $2, status = $3, stock = $4, 
				updated_at = NOW()
			WHERE id = $5`

		_, err = tx.ExecContext(ctx, updateMainQuery,
			product.SellingPrice,
			product.CostPrice,
			product.Status,
			product.Stock,
			existingProductID)

		if err != nil {
			log.Printf("UPDATE ERROR main product %s: %v", product.ProviderCode, err)
			return err
		}
	}

	return nil
}

func (pe *ProductExternal) getDenomination(product *digiflazz.InternalProduct) string {
	// Extract denomination dari product name
	// Contoh: "Telkomsel 10000" -> "10000"
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(product.ProviderName, -1)
	if len(matches) > 0 {
		return matches[0]
	}
	return "0"
}

func (pe *ProductExternal) getDenominationType(product *digiflazz.InternalProduct) string {
	name := strings.ToLower(product.ProviderName)

	if strings.Contains(name, "data") {
		return "data"
	} else if strings.Contains(name, "pln") {
		return "utility"
	} else if strings.Contains(name, "diamonds") {
		return "diamonds"
	} else if strings.Contains(name, "crystals") {
		return "crystals"
	} else if strings.Contains(name, "opals") {
		return "opals"
	} else if strings.Contains(name, "vouchers") {
		return "vouchers"
	} else if strings.Contains(name, "pulsa") || strings.Contains(name, "telkomsel") || strings.Contains(name, "xl") {
		return "pulsa"
	}

	return "other"
}

func (pe *ProductExternal) getSortOrder(product *digiflazz.InternalProduct) int {
	// Sort order berdasarkan harga atau kategori
	if product.SellingPrice <= 10000 {
		return 1
	} else if product.SellingPrice <= 50000 {
		return 2
	} else {
		return 3
	}
}
