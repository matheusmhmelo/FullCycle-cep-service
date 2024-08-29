package main

import (
	"context"
	"fmt"
	"github.com/matheusmhmelo/FullCycle-cep-api/internal/infra/gateway"
	"github.com/matheusmhmelo/FullCycle-cep-api/internal/infra/web"
	"github.com/matheusmhmelo/FullCycle-cep-api/internal/infra/web/webserver"
	"github.com/matheusmhmelo/FullCycle-cep-api/internal/usecase"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"log"
	"os"
)

func main() {
	viper.AutomaticEnv()

	shutdown, err := initProvider(viper.GetString("ZIPKIN_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	tracer := otel.Tracer("cep-api")

	gt := gateway.New(viper.GetString("API_KEY"), tracer)
	weatherUseCase := usecase.NewWeatherUseCase(gt)

	webServer := webserver.NewWebServer(viper.GetString("WEB_SERVER_PORT"))
	webOrderHandler := web.NewOrderHandler(weatherUseCase, tracer)
	webServer.AddHandler("/weather", webserver.HTTP_GET, webOrderHandler.Get)
	fmt.Println("Starting web server on port", viper.GetString("WEB_SERVER_PORT"))
	webServer.Start()
}

func initProvider(url string) (func(ctx context.Context) error, error) {
	logger := log.New(os.Stderr, "zipkin-cep-api", log.Ldate|log.Ltime|log.Llongfile)
	exporter, err := zipkin.New(
		url,
		zipkin.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("cep-api"),
		)),
	)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}
