package repository

import (
	"context"
	"database/sql"

	"github.com/wafi04/otomaxv2/internal/model"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}
func (pr *ProductRepository) GetProducts(ctx context.Context, limit, offset int) ([]*model.ProductData, error) {
    query := `
        SELECT 
            p.id,
            p.name AS product_name,
            p.description,
            p.price,
            p.original_price,
            p.status,
            p.image,
            pp.provider_id,
            pp.provider_code,
            pp.provider_name,
            pp.cost_price,
            pp.selling_price,
            pp.profit_margin,
            pp.stock,
            pp.is_available,
            pp.is_maintenance
        FROM products p
        JOIN provider_products pp ON pp.product_id = p.id
        ORDER BY p.id
        LIMIT $1 OFFSET $2
    `

    rows, err := pr.DB.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    productsMap := make(map[int]*model.ProductData)

    for rows.Next() {
        var (
            id            int
            name          string
            description   string
            price         float64
            originalPrice float64
            status        string
            image         *string
            provider      model.Provider
        )

        err := rows.Scan(
            &id,
            &name,
            &description,
            &price,
            &originalPrice,
            &status,
            &image,
            &provider.ProviderID,
            &provider.ProviderCode,
            &provider.ProviderName,
            &provider.CostPrice,
            &provider.SellingPrice,
            &provider.ProfitMargin,
            &provider.Stock,
            &provider.IsAvailable,
            &provider.IsMaintenance,
        )
        if err != nil {
            return nil, err
        }

        // Kalau product belum ada di map, buat dulu
        prod, exists := productsMap[id]
        if !exists {
            prod = &model.ProductData{
                ID:            id,
                Name:          name,
                Description:   description,
                Price:         price,
                OriginalPrice: originalPrice,
                Status:        status,
                Image:         image,
                Providers:     []model.Provider{},
            }
            productsMap[id] = prod
        }

        // Append provider ke slice
        prod.Providers = append(prod.Providers, provider)
    }

    // Ubah map ke slice
    products := make([]*model.ProductData, 0, len(productsMap))
    for _, p := range productsMap {
        products = append(products, p)
    }

    return products, nil
}
