package handler

import (
	"log"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/otomaxv2/internal/integrations/digiflazz"
	"github.com/wafi04/otomaxv2/internal/services/productexternal"
)

type ProductExternalHandler struct {
	productExternal *productexternal.ProductExternal
}

func NewProductExternalHandler(pe *productexternal.ProductExternal) *ProductExternalHandler {
	return &ProductExternalHandler{
		productExternal: pe,
	}
}

// Handler
func (peh *ProductExternalHandler) GetByDigiflazz(c *gin.Context) {
	// 1. Get raw data dari Digiflazz API
	digiflazzProducts, err := peh.productExternal.DigiflazzService.CheckPrice()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to get products: " + err.Error(),
		})
		return
	}

	// 2. Map ke format internal object
	mappedProducts := peh.mapDigiflazzToInternalProducts(digiflazzProducts)

	// 3. Save/process mapped products ke service
	processedProducts, err := peh.productExternal.GetProductDigiflazz(c, mappedProducts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to process products: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Products retrieved and processed successfully",
		"data":    processedProducts,
		"count":   len(processedProducts),
	})
}
func (peh *ProductExternalHandler) mapDigiflazzToInternalProducts(digiflazzProducts []*digiflazz.ProductData) []*digiflazz.InternalProduct {
	var internalProducts []*digiflazz.InternalProduct

	for _, dp := range digiflazzProducts {
		internal := &digiflazz.InternalProduct{
			ProviderCode: dp.BuyerSkuCode,
			ProviderName: dp.ProductName,
			Category:     dp.Category,
			Brand:        dp.Brand,
			Type:         dp.Type,
			Description:  dp.Desc,
			CostPrice:    dp.Price,
			Stock:        dp.Stock,
			IsUnlimited:  dp.UnlimitedStock,
			IsActive:     dp.BuyerProductStatus && dp.SellerProductStatus,
			SellerName:   dp.SellerName,
			StartCutOff:  dp.StartCutOff,
			EndCutOff:    dp.EndCutOff,
			SupportMulti: dp.Multi,
			Provider:     "digiflazz",
		}

		// HITUNG selling price dan profit margin DULU
		internal.SellingPrice = peh.calculateSellingPrice(internal.CostPrice)
		internal.ProfitMargin = peh.calculateProfitMargin(internal.CostPrice, internal.SellingPrice)
		internal.Status = peh.determineProductStatus(dp)

		// DEBUG SETELAH nilai sudah dihitung
		log.Printf("Debug: Product %s - Cost: %d, Selling: %d, Margin: %.2f%%",
			dp.BuyerSkuCode,
			internal.CostPrice,
			internal.SellingPrice,
			internal.ProfitMargin) // Atau gunakan rumus manual jika ProfitMargin masih int

		// Alternatif debug dengan manual calculation:
		if internal.CostPrice > 0 {
			manualMargin := float64(internal.SellingPrice-internal.CostPrice) * 100.0 / float64(internal.CostPrice)
			log.Printf("Manual margin calc for %s: %.2f%% (Raw: %.0f)",
				dp.BuyerSkuCode, manualMargin, math.Round(manualMargin))
		}

		// Check untuk potential overflow
		if internal.CostPrice > 2147483647 || internal.SellingPrice > 2147483647 {
			log.Printf("WARNING: Potential overflow for %s - Cost: %d, Selling: %d",
				dp.BuyerSkuCode, internal.CostPrice, internal.SellingPrice)
		}

		// Check profit margin calculation
		marginFloat := float64(internal.SellingPrice-internal.CostPrice) * 100.0 / float64(internal.CostPrice)
		if math.Abs(marginFloat) > 2147483647 {
			log.Printf("WARNING: ProfitMargin overflow for %s - Calculated: %.2f",
				dp.BuyerSkuCode, marginFloat)
		}

		internalProducts = append(internalProducts, internal)
	}

	return internalProducts
}

// Tambahan: Debug function untuk test specific products
func (peh *ProductExternalHandler) debugSpecificProduct(productCode string, costPrice int) {
	sellingPrice := peh.calculateSellingPrice(costPrice)
	profitMargin := peh.calculateProfitMargin(costPrice, sellingPrice)

	log.Printf("=== DEBUG %s ===", productCode)
	log.Printf("Cost Price: %d", costPrice)
	log.Printf("Selling Price: %d", sellingPrice)
	log.Printf("Profit Margin: %v", profitMargin)

	// Manual calculation check
	if costPrice > 0 {
		manualMargin := float64(sellingPrice-costPrice) * 100.0 / float64(costPrice)
		log.Printf("Manual Margin: %.2f%%", manualMargin)
		log.Printf("Manual Rounded: %.0f", math.Round(manualMargin))
	}
	log.Printf("==================")
}

// Test untuk products yang bermasalah
func (peh *ProductExternalHandler) testProblematicProducts() {
	// Test dengan products yang error dari log
	testProducts := []struct {
		code  string
		price int
	}{
		{"ax10", 10000},    // Contoh price
		{"ff12", 15000},    // Contoh price
		{"ml10", 20000},    // Contoh price
		{"pln100", 100000}, // Contoh price
	}

	for _, tp := range testProducts {
		peh.debugSpecificProduct(tp.code, tp.price)
	}
}
func (peh *ProductExternalHandler) calculateSellingPrice(costPrice int) int {
	// Contoh: markup 15%
	return int(math.Round(float64(costPrice) * 1.15))
}

func (peh *ProductExternalHandler) calculateProfitMargin(costPrice, sellingPrice int) int {
	if costPrice == 0 {
		return 0
	}
	return int(math.Round(float64(sellingPrice-costPrice) * 100 / float64(costPrice)))
}

func (peh *ProductExternalHandler) determineProductStatus(dp *digiflazz.ProductData) string {
	if !dp.BuyerProductStatus || !dp.SellerProductStatus {
		return "inactive"
	}
	if dp.Stock == 0 && !dp.UnlimitedStock {
		return "out_of_stock"
	}
	return "active"
}
