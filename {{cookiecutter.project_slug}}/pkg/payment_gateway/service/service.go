package payment_gateway

import (
	"{{ cookiecutter.module_path }}/pkg/payment_gateway"
	c "{{ cookiecutter.module_path }}/utils"
)

func VerifyPayment(PaymentRef string, PaymentProvider string,cfg c.Config) (string, any,error) {

	payment, err := payment_gateway.NewPayment(PaymentProvider, cfg)
	if err != nil {
		return "", nil, err
	}

	res, resData, err := payment.VerifyPayment(PaymentRef)
	return res, resData,err
}
