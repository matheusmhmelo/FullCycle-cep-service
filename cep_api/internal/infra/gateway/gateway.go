package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"net/url"
)

var (
	ErrorInvalidCEP  = errors.New("invalid zipcode")
	ErrorNotFoundCEP = errors.New("can not find zipcode")
)

var doFunc = http.DefaultClient.Do

type WeatherGatewayInterface interface {
	ValidateLocation(cep string) (string, error)
	GetWeather() (float64, error)
}

type cepResponse struct {
	Error    string `json:"erro"`
	Location string `json:"localidade"`
}

type weatherResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type weatherGatewayImpl struct {
	apiKey     string
	location   string
	otelTracer trace.Tracer
}

func New(apiKey string, otelTracer trace.Tracer) WeatherGatewayInterface {
	return &weatherGatewayImpl{
		apiKey:     apiKey,
		otelTracer: otelTracer,
	}
}

func (w *weatherGatewayImpl) ValidateLocation(cep string) (string, error) {
	ctx, span := w.otelTracer.Start(context.Background(), "Location Request")

	if len(cep) != 8 {
		return "", ErrorInvalidCEP
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "https://viacep.com.br/ws/"+cep+"/json/", nil)
	if err != nil {
		return "", err
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := doFunc(req)
	span.End()
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusBadRequest {
		return "", ErrorNotFoundCEP
	}

	var content cepResponse
	err = json.NewDecoder(resp.Body).Decode(&content)
	if err != nil {
		return "", err
	}

	if content.Error != "" {
		return "", ErrorNotFoundCEP
	}
	w.location = content.Location
	return content.Location, nil
}

func (w *weatherGatewayImpl) GetWeather() (float64, error) {
	ctx, span := w.otelTracer.Start(context.Background(), "Weather Request")

	u, err := url.Parse("http://api.weatherapi.com/v1/current.json")
	if err != nil {
		return 0, err
	}

	query := url.Values{}
	query.Set("key", w.apiKey)
	query.Set("q", w.location)
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return 0, err
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := doFunc(req)
	span.End()
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("invalid status received")
	}

	var content weatherResponse
	err = json.NewDecoder(resp.Body).Decode(&content)
	if err != nil {
		return 0, err
	}
	return content.Current.TempC, nil
}
