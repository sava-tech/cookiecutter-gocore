package payment_gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type FlutterWave struct {
	SecretKey string
	BaseURL   string
}

func NewFlutterWave() *FlutterWave {
	secret := os.Getenv("FLUTTERWAVE_SECRET_KEY")
	testSecret := os.Getenv("FLUTTERWAVE_SECRET_KEY_TEST")
	testPayment := os.Getenv("TEST_PAYMENT") // "true" or "false"

	if testPayment == "true" {
		secret = testSecret
	}

	return &FlutterWave{
		SecretKey: secret,
		BaseURL:   "https://api.flutterwave.com/v3/transactions/",
	}
}

type FlutterWaveVerifyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (fw *FlutterWave) VerifyPayment(paymentRef string) (string, any, error) {
	path := fmt.Sprintf("verify_by_reference?tx_ref=%s", paymentRef)
	url := fw.BaseURL + path

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nil, err
	}

	req.Header.Set("Authorization", "Bearer "+fw.SecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	var response FlutterWaveVerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", nil, err
	}

	// (Optional) Here is where you'd send an email → use your own email service
	// sendEmail("Response", fmt.Sprintf("%v", response))

	if resp.StatusCode == 200 {
		return response.Status, response.Data, nil
	}

	if resp.StatusCode != 200 {
		log.Println("Error verifying payment:", response.Message)
		return response.Status, response.Message, fmt.Errorf("error verifying payment: %s", response.Message)
	}

	return response.Status, response.Data, nil
}
