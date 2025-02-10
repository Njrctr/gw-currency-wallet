package exchanger_grpc

import (
	"context"
	"fmt"

	exchangev1 "github.com/Njrctr/gw-proto-exchange/gen/go/exchange"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	api exchangev1.ExchangeServiceClient
}

func NewGRPCClient(
	ctx context.Context,
	addr string,
) (*GRPCClient, error) {

	con, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc.NewGRPCClient: %W", err)
	}

	return &GRPCClient{
		api: exchangev1.NewExchangeServiceClient(con),
	}, nil
}

func (c *GRPCClient) GetExchangeRates(ctx context.Context) (exchangev1.ExchangeRatesResponse, error) {
	rates, err := c.api.GetExchangeRates(ctx, &exchangev1.Empty{})
	if err != nil {
		return exchangev1.ExchangeRatesResponse{}, fmt.Errorf("grpc.GetExchangeRates: %w", err)
	}

	return exchangev1.ExchangeRatesResponse{
		Rates: rates.Rates,
	}, nil
}

func (c *GRPCClient) GetRateForCurrency(ctx context.Context, from, to string) (exchangev1.ExchangeRateResponse, error) {
	rate, err := c.api.GetExchangeRateForCurrency(ctx, &exchangev1.CurrencyRequest{
		FromCurrency: from,
		ToCurrency:   to,
	})
	if err != nil {
		return exchangev1.ExchangeRateResponse{}, fmt.Errorf("grpc.GetExchangeRateForCurrency: %w", err)
	}

	return exchangev1.ExchangeRateResponse{
		FromCurrency: rate.FromCurrency,
		ToCurrency:   rate.ToCurrency,
		Rate:         rate.Rate,
	}, nil
}
