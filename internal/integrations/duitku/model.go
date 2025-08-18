package duitku

import "net/http"

type DuitkuCreateTransactionParams struct {
	PaymentAmount   int     `json:"paymentAmount"`
	MerchantOrderId string  `json:"merchantOrderId"`
	ProductDetails  string  `json:"productDetails"`
	PaymentCode     string  `json:"paymentCode"`
	Cust            *string `json:"cust,omitempty"`
	CallbackUrl     *string `json:"callbackUrl,omitempty"`
	ReturnUrl       *string `json:"returnUrl,omitempty"`
}

type ResponseFromDuitkuCheckTransaction struct {
	Status int `json:"status"`
	Data   struct {
		MerchantOrderId string `json:"merchantOrderId"`
		Reference       string `json:"reference"`
		Amount          string `json:"amount"`
		Fee             string `json:"fee"`
		StatusCode      string `json:"statusCode"`
		StatusMessage   string `json:"statusMessage"`
	} `json:"data"`
}

type DuitkuCreateTransactionResponse struct {
	Status        string `json:"status"`
	Code          string `json:"code"`
	Message       string `json:"message"`
	QrString      string `json:"qrString,omitempty"`
	VANumber      string `json:"vaNumber,omitempty"`
	PaymentUrl    string `json:"paymentUrl,omitempty"`
	Amount        string `json:"amount,omitempty"`
	Reference     string `json:"reference,omitempty"`
	StatusCode    string `json:"statusCode,omitempty"`
	StatusMessage string `json:"statusMessage,omitempty"`
}

type DuitkuService struct {
	DuitkuKey             string
	DuitkuMerchantCode    string
	DuitkuExpiryPeriod    *int64
	BaseUrl               string
	SandboxUrl            string
	BaseUrlGetTransaction string
	BaseUrlGetBalance     string
	HttpClient            *http.Client
}

type PaymentResponse struct {
	MerchantOrderId string `json:"merchantOrderId"`
	Signature       string `json:"signature"`
	Timestamp       string `json:"timestamp"`
	PaymentUrl      string `json:"paymentUrl"`
	QrString        string `json:"qrString,omitempty"`
	VANumber        string `json:"vaNumber,omitempty"`
	Amount          string `json:"amount"`
	Reference       string `json:"reference"`
	StatusCode      string `json:"statusCode"`
	StatusMessage   string `json:"statusMessage"`
}
