package app

import (
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/charge"
	"os"
)

var (
	stripeKey     = os.Getenv("STRIPE_PRIVATE_KEY")
	mailgunDomain = os.Getenv("MAILGUN_DOMAIN")
	mailgunKey    = os.Getenv("MAILGUN_PRIVATE_KEY")
)

type Activities struct {}

func (a *Activities) CreateStripeCharge(_ context.Context, cart CartState) error {
	stripe.Key = stripeKey
	var amount float32 = 0
	var description string = ""
	for _, item := range cart.Items {
		var product Product
		for _, _product := range Products {
			if _product.Id == item.ProductId {
				product = _product
				break
			}
		}
		amount += float32(item.Quantity) * product.Price
		if len(description) > 0 {
			description += ", "
		}
		description += product.Name
	}

	_, err := charge.New(&stripe.ChargeParams{
		Amount:       stripe.Int64(int64(amount * 100)),
		Currency:     stripe.String(string(stripe.CurrencyUSD)),
		Description:  stripe.String(description),
		Source:       &stripe.SourceParams{Token: stripe.String("tok_visa")},
		ReceiptEmail: stripe.String(cart.Email),
	})

	if err != nil {
		fmt.Println("Stripe err: " + err.Error())
	}

	return err
}

func (a *Activities) SendAbandonedCartEmail(_ context.Context, email string) error {
	mg := mailgun.NewMailgun(mailgunDomain, mailgunKey)
	m := mg.NewMessage(
		"noreply@"+mailgunDomain,
		"You've abandoned your shopping cart!",
		"Go to http://localhost:8080 to finish checking out!",
		email,
	)
	_, _, err := mg.Send(m)
	if err != nil {
		fmt.Println("Mailgun err: " + err.Error())
		return err
	}

	return err
}
