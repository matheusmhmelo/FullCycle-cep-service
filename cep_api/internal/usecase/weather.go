package usecase

import (
	"context"
	"github.com/matheusmhmelo/FullCycle-cep-api/internal/infra/gateway"
)

type Weather struct {
	City       string  `json:"city"`
	Fahrenheit float64 `json:"temp_F"`
	Celsius    float64 `json:"temp_C"`
	Kelvin     float64 `json:"temp_k"`
}

type WeatherUseCase interface {
	Execute(ctx context.Context, cep string) (*Weather, error)
}

type weatherUseCaseImpl struct {
	Gateway gateway.WeatherGatewayInterface
}

func NewWeatherUseCase(
	Gateway gateway.WeatherGatewayInterface,
) WeatherUseCase {
	return &weatherUseCaseImpl{
		Gateway: Gateway,
	}
}

func (w *weatherUseCaseImpl) Execute(ctx context.Context, cep string) (*Weather, error) {
	loc, err := w.Gateway.ValidateLocation(ctx, cep)
	if err != nil {
		return nil, err
	}

	weatherC, err := w.Gateway.GetWeather(ctx)
	if err != nil {
		return nil, err
	}

	return &Weather{
		City:       loc,
		Fahrenheit: (weatherC * 1.8) + 32,
		Celsius:    weatherC,
		Kelvin:     weatherC + 273,
	}, nil
}
