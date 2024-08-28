package web

import (
	"encoding/json"
	"errors"
	"github.com/matheusmhmelo/FullCycle-cep-service/cep_api/internal/infra/gateway"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"

	"github.com/matheusmhmelo/FullCycle-cep-service/cep_api/internal/usecase"
)

type OrderHandler struct {
	weather    usecase.WeatherUseCase
	otelTracer trace.Tracer
}

func NewOrderHandler(
	weather usecase.WeatherUseCase,
	otelTracer trace.Tracer,
) *OrderHandler {
	return &OrderHandler{
		weather:    weather,
		otelTracer: otelTracer,
	}
}

func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := h.otelTracer.Start(ctx, "POST")
	defer span.End()

	cep := r.URL.Query().Get("cep")
	output, err := h.weather.Execute(cep)
	if err != nil {
		if errors.Is(err, gateway.ErrorInvalidCEP) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, gateway.ErrorNotFoundCEP) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
