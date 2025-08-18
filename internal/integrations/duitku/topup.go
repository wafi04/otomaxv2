package duitku

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/wafi04/otomaxv2/internal/config"
)

func NewDuitkuService(cfg *config.Config) *DuitkuService {
	return &DuitkuService{
		DuitkuKey:             cfg.PaymentGateway.DuitkuConfig.DuitkuKey,
		DuitkuMerchantCode:    cfg.PaymentGateway.DuitkuConfig.DuitkuMerchantCode,
		BaseUrl:               "https://passport.duitku.com/webapi/api/merchant/v2/inquiry",
		SandboxUrl:            "https://sandbox.duitku.com/webapi/api/merchant/v2/inquiry",
		BaseUrlGetTransaction: "https://passport.duitku.com/webapi/api/merchant/transactionStatus",
		BaseUrlGetBalance:     "https://passport.duitku.com/webapi/api/disbursement/checkbalance",
		HttpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *DuitkuService) CreateTransaction(ctx context.Context, params *DuitkuCreateTransactionParams) (*DuitkuCreateTransactionResponse, error) {

	signature := s.generateSignature(params.MerchantOrderId, params.PaymentAmount)

	payload := map[string]interface{}{
		"merchantCode":    s.DuitkuMerchantCode,
		"paymentAmount":   params.PaymentAmount,
		"merchantOrderId": params.MerchantOrderId,
		"productDetails":  params.ProductDetails,
		"paymentMethod":   params.PaymentCode,
		"signature":       signature,
		"callbackUrl":     params.CallbackUrl,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, nil
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.SandboxUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, nil
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil
	}

	var duitkuResponse DuitkuCreateTransactionResponse
	if err := json.Unmarshal(body, &duitkuResponse); err != nil {
		return s.createErrorResponse(params.MerchantOrderId, fmt.Sprintf("Failed to parse success response: %v", err)), nil
	}
	return &duitkuResponse, nil
}

func (s *DuitkuService) generateSignature(merchantOrderId string, paymentAmount int) string {

	signatureString := s.DuitkuMerchantCode + merchantOrderId + strconv.Itoa(paymentAmount) + s.DuitkuKey

	h := md5.New()
	h.Write([]byte(signatureString))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *DuitkuService) createErrorResponse(merchantOrderId, errorMessage string) *DuitkuCreateTransactionResponse {
	return &DuitkuCreateTransactionResponse{
		Status:  "false",
		Message: errorMessage,
	}
}
