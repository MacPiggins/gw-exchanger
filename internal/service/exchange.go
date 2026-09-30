package service

import "context"

type ExchangeRateStorage interface {
	GetExchangeRate(ctx context.Context, from, to string) (float32, error)
	GetExchangeRates(ctx context.Context) (map[string]float32, error)
}

type ExchangeService struct {
	storage ExchangeRateStorage
}

func NewExchangeService(storage ExchangeRateStorage) *ExchangeService {
	return &ExchangeService{storage: storage}
}

func (svc *ExchangeService) GetExchangeRate(ctx context.Context, from, to string) (float32, error) {
	return svc.storage.GetExchangeRate(ctx, from, to)
}
func (svc *ExchangeService) GetExchangeRates(ctx context.Context) (map[string]float32, error) {
	return svc.storage.GetExchangeRates(ctx)
}
