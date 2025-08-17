package digiflazz

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func NewDigiflazzService(config DigiConfig) *DigiflazzService {
	return &DigiflazzService{
		config: config,
	}
}

func (d *DigiflazzService) generateSign(username, apiKey, cmd string) string {
	data := username + apiKey + cmd
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func (d *DigiflazzService) CheckPrice() ([]*ProductData, error) {
	sign := d.generateSign(d.config.DigiUsername, d.config.DigiKey, "pricelist")

	requestPayload := map[string]interface{}{
		"username": d.config.DigiUsername,
		"cmd":      "pricelist",
		"sign":     sign,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Make HTTP request
	resp, err := http.Post("https://api.digiflazz.com/v1/price-list", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var apiResponse DigiflazzResponse
	if err := json.Unmarshal(body, &apiResponse); err == nil && len(apiResponse.Data) > 0 {
		result := make([]*ProductData, len(apiResponse.Data))
		for i := range apiResponse.Data {
			result[i] = &apiResponse.Data[i]
		}
		return result, nil
	}

	var directArray []ProductData
	if err := json.Unmarshal(body, &directArray); err == nil {
		result := make([]*ProductData, len(directArray))
		for i := range directArray {
			result[i] = &directArray[i]
		}
		return result, nil
	}

	return nil, fmt.Errorf("failed to unmarshal response. Response body: %s", string(body))
}
