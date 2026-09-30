package service

import (
	"context"
	"errors"
	"testing"
)

type exchangeRateStorageMock struct {
	getExchangeRate  func(context.Context, string, string) (float32, error)
	getExchangeRates func(context.Context) (map[string]float32, error)
}

func (m *exchangeRateStorageMock) GetExchangeRate(ctx context.Context, from, to string) (float32, error) {
	return m.getExchangeRate(ctx, from, to)
}

func (m *exchangeRateStorageMock) GetExchangeRates(ctx context.Context) (map[string]float32, error) {
	return m.getExchangeRates(ctx)
}

func TestExchangeServiceGetExchangeRate(t *testing.T) {
	ctx := context.Background()
	wantRate := float32(1.25)
	var gotCtx context.Context
	var gotFrom, gotTo string

	storage := &exchangeRateStorageMock{
		getExchangeRate: func(ctx context.Context, from, to string) (float32, error) {
			gotCtx, gotFrom, gotTo = ctx, from, to
			return wantRate, nil
		},
	}

	gotRate, err := NewExchangeService(storage).GetExchangeRate(ctx, "USD", "EUR")
	if err != nil {
		t.Fatalf("GetExchangeRate() error = %v", err)
	}
	if gotRate != wantRate {
		t.Errorf("GetExchangeRate() = %v, want %v", gotRate, wantRate)
	}
	if gotCtx != ctx || gotFrom != "USD" || gotTo != "EUR" {
		t.Errorf("storage called with (%v, %q, %q), want (%v, %q, %q)", gotCtx, gotFrom, gotTo, ctx, "USD", "EUR")
	}
}

func TestExchangeServiceGetExchangeRatePropagatesError(t *testing.T) {
	wantErr := errors.New("storage failure")
	storage := &exchangeRateStorageMock{
		getExchangeRate: func(context.Context, string, string) (float32, error) {
			return 0, wantErr
		},
	}

	_, err := NewExchangeService(storage).GetExchangeRate(context.Background(), "USD", "EUR")
	if !errors.Is(err, wantErr) {
		t.Errorf("GetExchangeRate() error = %v, want %v", err, wantErr)
	}
}

func TestExchangeServiceGetExchangeRates(t *testing.T) {
	ctx := context.Background()
	wantRates := map[string]float32{"USD/EUR": 0.92, "EUR/GBP": 0.86}
	var gotCtx context.Context
	storage := &exchangeRateStorageMock{
		getExchangeRates: func(ctx context.Context) (map[string]float32, error) {
			gotCtx = ctx
			return wantRates, nil
		},
	}

	gotRates, err := NewExchangeService(storage).GetExchangeRates(ctx)
	if err != nil {
		t.Fatalf("GetExchangeRates() error = %v", err)
	}
	if gotCtx != ctx {
		t.Errorf("storage called with context %v, want %v", gotCtx, ctx)
	}
	if len(gotRates) != len(wantRates) {
		t.Fatalf("GetExchangeRates() = %v, want %v", gotRates, wantRates)
	}
	for key, want := range wantRates {
		if gotRates[key] != want {
			t.Errorf("GetExchangeRates()[%q] = %v, want %v", key, gotRates[key], want)
		}
	}
}

func TestExchangeServiceGetExchangeRatesPropagatesError(t *testing.T) {
	wantErr := errors.New("storage failure")
	storage := &exchangeRateStorageMock{
		getExchangeRates: func(context.Context) (map[string]float32, error) {
			return nil, wantErr
		},
	}

	_, err := NewExchangeService(storage).GetExchangeRates(context.Background())
	if !errors.Is(err, wantErr) {
		t.Errorf("GetExchangeRates() error = %v, want %v", err, wantErr)
	}
}
