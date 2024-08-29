package usecase

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"io"
	"net/http"
	"net/url"
)

var (
	doFunc = http.DefaultClient.Do
)

type Weather struct {
	City       string  `json:"city"`
	Fahrenheit float64 `json:"temp_F"`
	Celsius    float64 `json:"temp_C"`
	Kelvin     float64 `json:"temp_k"`
}

type WeatherUseCase interface {
	Execute(ctx context.Context, cep string) ([]byte, int, error)
}

type weatherUseCaseImpl struct {
	baseURL string
}

func NewWeatherUseCase(baseURL string) WeatherUseCase {
	return &weatherUseCaseImpl{
		baseURL: baseURL,
	}
}

func (w *weatherUseCaseImpl) Execute(ctx context.Context, cep string) ([]byte, int, error) {
	u, err := url.Parse(w.baseURL)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	query := url.Values{}
	query.Set("cep", cep)
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := doFunc(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return bodyBytes, resp.StatusCode, nil
}
