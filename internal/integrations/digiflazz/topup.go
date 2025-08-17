package digiflazz

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func (d *DigiflazzService) TopUp(ctx context.Context, req CreateTransactionToDigiflazz) (*TransactionCreateDigiflazzResponse, error) {
	log.Print("Transaction called digiflazz")
	data := d.config.DigiUsername + d.config.DigiKey + req.RefID
	hash := md5.Sum([]byte(data))
	sign := fmt.Sprintf("%x", hash)

	requestPayload := map[string]interface{}{
		"username":       d.config.DigiUsername,
		"buyer_sku_code": req.BuyerSKUCode,
		"customer_no":    req.CustomerNo,
		"ref_id":         req.RefID,
		"sign":           sign,
		"cb_url":         "https://e3abda9a4f73.ngrok-free.app/api/transactions/callback/digiflazz",
	}

	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.digiflazz.com/v1/transaction", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", "DigiflazzClient/1.0")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResponse TransactionCreateDigiflazzResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, body: %s", err, string(body))
	}

	switch apiResponse.Data.Status {
	case "Sukses":
		return &apiResponse, nil
	case "Pending":
		return &apiResponse, nil
	case "Gagal":
		return &apiResponse, nil
	default:
		return &apiResponse, nil
	}

}
