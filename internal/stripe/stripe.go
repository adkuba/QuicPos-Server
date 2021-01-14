package stripe

import (
	"QuicPos/internal/data"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
)

var stripeKey = data.StripePrivate

//Init stripe
func Init() {
	stripe.Key = stripeKey
}

//CreatePaymentIntent creates client secret key for specified amount
func CreatePaymentIntent(amount float64) (string, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount * 100)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return "", err
	}

	return pi.ClientSecret, nil
}
