package digiflazz

type ProductData struct {
	BuyerProductStatus  bool   `json:"buyer_product_status"`
	BuyerSkuCode        string `json:"buyer_sku_code"`
	Category            string `json:"category"`
	Desc                string `json:"desc"`
	EndCutOff           string `json:"end_cut_off"`
	Multi               bool   `json:"multi"`
	Price               int    `json:"price"`
	ProductName         string `json:"product_name"`
	SellerName          string `json:"seller_name"`
	SellerProductStatus bool   `json:"seller_product_status"`
	StartCutOff         string `json:"start_cut_off"`
	Stock               int    `json:"stock"`
	Type                string `json:"type"`
	UnlimitedStock      bool   `json:"unlimited_stock"`
	Brand               string `json:"brand"`
}

type CreateTransactionToDigiflazz struct {
	BuyerSKUCode string `json:"buyer_sku_code"`
	CustomerNo   string `json:"customer_no"`
	RefID        string `json:"ref_id"`
	CallbackURL  string `json:"cb_url,omitempty"`
}
type TransactionCreateDigiflazzResponse struct {
	Data struct {
		RefID          string `json:"ref_id"`
		CustomerNo     string `json:"customer_no"`
		BuyerSKUCode   string `json:"buyer_sku_code"`
		Message        string `json:"message"`
		Status         string `json:"status"`
		RC             string `json:"rc"`
		SN             string `json:"sn"`
		BuyerLastSaldo int    `json:"buyer_last_saldo"`
		Price          int    `json:"price"`
		Tele           string `json:"tele"`
		WA             string `json:"wa"`
	} `json:"data"`
}

type DigiflazzErrorResponse struct {
	Data struct {
		Message string `json:"message"`
		RC      string `json:"rc"`
	} `json:"data"`
}

type DigiflazzResponse struct {
	Data []ProductData `json:"data"`
}

type DigiConfig struct {
	DigiKey      string
	DigiUsername string
}

type DigiflazzService struct {
	config DigiConfig
}

type InternalProduct struct {
	ProviderCode  string `json:"provider_code"`
	ProviderName  string `json:"provider_name"`
	Category      string `json:"category"`
	Brand         string `json:"brand"`
	Type          string `json:"type"`
	Description   string `json:"description"`
	CostPrice     int    `json:"cost_price"`
	SellingPrice  int    `json:"selling_price"`
	ProfitMargin  int    `json:"profit_margin"`
	Stock         int    `json:"stock"`
	IsUnlimited   bool   `json:"is_unlimited"`
	IsActive      bool   `json:"is_active"`
	Status        string `json:"status"`
	SellerName    string `json:"seller_name"`
	StartCutOff   string `json:"start_cut_off"`
	EndCutOff     string `json:"end_cut_off"`
	SupportMulti  bool   `json:"support_multi"`
	Provider      string `json:"provider"`
	CategoryID    *int   `json:"categoryId"`
	SubCategoryID *int   `json:"subCategoryId"`
}
