package model

type Category struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	SubName         string  `json:"subName"`
	Brand           string  `json:"brand"`
	Code            string  `json:"code"`
	IsCheckNickname string  `json:"isCheckNickname"`
	Status          string  `json:"status"`
	Thumbnail       string  `json:"thumbnail"`
	Type            string  `json:"type"`
	Banner          string  `json:"banner"`
	Instruction     *string `json:"instruction,omitempty"`
	Information     *string `json:"information,omitempty"`
	Placeholder1    string  `json:"placeholder1"`
	Placeholder2    *string `json:"placeholder2,omitempty"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

type Product struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Price            int     `json:"price"`
	DenominationType string  `json:"denominationType"`
	SubCategoryName  *string `json:"subCategoryName,omitempty"`
	SubCategoryID    *int    `json:"subCategoryID,omitempty"`
}
type CategoryCodeResponse struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	SubName         string    `json:"subName"`
	Brand           string    `json:"brand"`
	Code            string    `json:"code"`
	IsCheckNickname string    `json:"isCheckNickname"`
	Status          string    `json:"status"`
	Thumbnail       string    `json:"thumbnail"`
	Type            string    `json:"type"`
	Banner          string    `json:"banner"`
	Instruction     *string   `json:"instruction,omitempty"`
	Information     *string   `json:"information,omitempty"`
	Products        []Product `json:"products"`
	SubCategories   []SubCategory
}

type CreateCategory struct {
	Name            string  `json:"name"`
	SubName         string  `json:"subName"`
	Brand           string  `json:"brand"`
	Code            string  `json:"code"`
	IsCheckNickname string  `json:"isCheckNickname"`
	Status          string  `json:"status"`
	Thumbnail       string  `json:"thumbnail"`
	Type            string  `json:"type"`
	Banner          string  `json:"banner"`
	Placeholder1    string  `json:"placeholder1"`
	Placeholder2    *string `json:"placeholder2,omitempty"`
	Instruction     *string `json:"instruction,omitempty"`
	Information     string  `json:"information"`
}

type UpdateCategory struct {
	Name            *string `json:"name"`
	SubName         *string `json:"sub_name"`
	Brand           *string `json:"brand"`
	Code            *string `json:"code"`
	IsCheckNickname *string `json:"is_check_nickname"`
	Status          *string `json:"status"`
	Thumbnail       *string `json:"thumbnail"`
	Type            *string `json:"type"`
	Instruction     *string `json:"instruction"`
	Information     *string `json:"information"`
}
