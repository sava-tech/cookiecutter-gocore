package payment_gateway

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Paystack struct {
	SecretKey string
	BaseURL   string
}

func NewPaystack() *Paystack {
	secret := os.Getenv("PAYSTACK_SECRET_KEY")

	return &Paystack{
		SecretKey: secret,
		BaseURL:   "https://api.paystack.co",
	}
}

type PaystackVerifyResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (ps *Paystack) VerifyPayment(paymentRef string) (string, any, error) {
	path := fmt.Sprintf("/transaction/verify/%s", paymentRef)
	url := ps.BaseURL + path

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nil, err
	}

	req.Header.Set("Authorization", "Bearer "+ps.SecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	var response PaystackVerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", nil, err
	}

	if resp.StatusCode == 200 {
		return fmt.Sprint(response.Status), response.Data, nil
	}

	return fmt.Sprint(response.Status), response.Message, nil
}
