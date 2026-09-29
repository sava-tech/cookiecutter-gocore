package payment_gateway

import (
	"fmt"

	c "{{ cookiecutter.module_path }}/utils"
)

type Config struct {
	PaymentProvider      string
	TestPayment          bool
	FlutterWaveSecretKey string
	FlutterWaveBaseURL   string
	PaystackBaseURL      string
	PaystackSecretKey    string
}

func NewPayment(PaymentProvider string, cfg c.Config) (Payment, error) {
	switch PaymentProvider {
	case "flutterwave":
		return &FlutterWave{SecretKey: cfg.FlutterWaveSecretKey, BaseURL: cfg.FlutterWaveBaseURL}, nil

	case "paystack":
		return &Paystack{SecretKey: cfg.PaystackSecretKey, BaseURL: cfg.PaystackBaseURL}, nil

	default:
		return nil, fmt.Errorf("unsupported payment provider: %s", cfg.Provider)
	}
}
