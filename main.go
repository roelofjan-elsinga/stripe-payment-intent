package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/roelofjan-elsinga/lambda/models"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"net/http"
	"os"
)

type LambdaRequest struct {
	Amount   int64           `json:"amount_cents"`
	Currency stripe.Currency `json:"currency"`
}

type LambdaResponse struct {
	ClientSecret string `json:"client_secret"`
}

func main() {
	lambda.Start(HttpHandler)
}

func HttpHandler(ctx context.Context, event models.GatewayRequest) (models.GatewayResponse, error) {

	var body LambdaRequest

	err := json.Unmarshal([]byte(event.Body), &body)

	if err != nil {
		return models.GatewayResponse{
			IsBase64Encoded: false,
			StatusCode:      http.StatusInternalServerError,
			Body:            err.Error(),
		}, err
	}

	responseBody, err := Action(ctx, body)

	if err != nil {
		return models.GatewayResponse{
			IsBase64Encoded: false,
			StatusCode:      http.StatusInternalServerError,
			Body:            err.Error(),
		}, err
	}

	jsonBody, err := json.Marshal(responseBody)

	if err != nil {
		return models.GatewayResponse{
			IsBase64Encoded: false,
			StatusCode:      http.StatusInternalServerError,
			Body:            err.Error(),
		}, err
	}

	return models.GatewayResponse{
		IsBase64Encoded: false,
		StatusCode:      http.StatusOK,
		Body:            string(jsonBody),
	}, nil
}

func Action(ctx context.Context, event LambdaRequest) (*LambdaResponse, error) {
	stripe.Key = os.Getenv("STRIPE_SECRET")

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(event.Amount),
		Currency: stripe.String(string(event.Currency)),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"ideal",
		}),
	}
	intent, err := paymentintent.New(params)

	if err != nil {
		return nil, err
	}

	return &LambdaResponse{
		ClientSecret: intent.ClientSecret,
	}, nil
}
