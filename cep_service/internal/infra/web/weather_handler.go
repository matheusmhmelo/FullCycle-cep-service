package web

import (
	"encoding/json"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"

	"github.com/matheusmhmelo/FullCycle-cep-service/internal/usecase"
)

type requestBody struct {
	CEP interface{} `json:"cep"`
}

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

func (h *OrderHandler) Post(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := h.otelTracer.Start(ctx, "POST")
	defer span.End()

	var body requestBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)

	cep, ok := body.CEP.(string)
	if !ok || len(cep) > 8 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	output, code, err := h.weather.Execute(ctx, cep)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	w.WriteHeader(code)
	_, err = w.Write(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
