// Interface definition

package payment_gateway

type Payment interface {
	// VerifyPayment : verifies payment using the ref
	VerifyPayment(paymentRef string) (string, any, error)
}
