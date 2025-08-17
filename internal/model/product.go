package model

type ProductData struct {
    ID            int  `json:"id"`
    Name          string `json:"name"`
    Description   string `json:"description"`
    Price         float64 `json:"price"`
    OriginalPrice float64 `json:"originalPrice"`
    Status        string `json:"status"`
    Image         *string `json:"image"`
    Providers     []Provider `json:"providers"`
}

type Provider struct {
    ProviderID    int
    ProviderCode  string  `json:"providerCode"`
    ProviderName  string `json:"providerName"`
    CostPrice     float64 `json:"costPrice"`
    SellingPrice  float64 `json:"sellingPrice"`
    ProfitMargin  float64 `json:"marginProfit"`
    Stock         int `json:"stock"`
    IsAvailable   bool `json:"isAvailable"`
    IsMaintenance bool `json:"isMaintenance"`
}
