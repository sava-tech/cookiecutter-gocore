package payment_gateway

import (
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/pkg/payment_gateway"
	c "github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/utils"
)

func VerifyPayment(PaymentRef string, PaymentProvider string,cfg c.Config) (string, any,error) {

	payment, err := payment_gateway.NewPayment(PaymentProvider, cfg)
	if err != nil {
		return "", nil, err
	}

	res, resData, err := payment.VerifyPayment(PaymentRef)
	return res, resData,err
}
