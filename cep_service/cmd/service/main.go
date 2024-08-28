package main

import (
	"context"
	"fmt"
	"github.com/matheusmhmelo/FullCycle-cep-service/cep_service/internal/infra/web"
	"github.com/matheusmhmelo/FullCycle-cep-service/cep_service/internal/infra/web/webserver"
	"github.com/matheusmhmelo/FullCycle-cep-service/cep_service/internal/usecase"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"os/signal"
	"time"
)

func init() {
	viper.AutomaticEnv()
}

func main() {
	signCh := make(chan os.Signal, 1)
	signal.Notify(signCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := initProvider(viper.GetString("OTEL_SERVICE_NAME"), viper.GetString("OTEL_EXPORTER_ODPL_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	tracer := otel.Tracer("cep-service")

	weatherUseCase := usecase.NewWeatherUseCase(viper.GetString("BASE_URL"))

	webServer := webserver.NewWebServer(viper.GetString("WEB_SERVER_PORT"))
	webOrderHandler := web.NewOrderHandler(weatherUseCase, tracer)
	webServer.AddHandler("/", webserver.HTTP_POST, webOrderHandler.Post)

	go func() {
		fmt.Println("Starting web server on port", viper.GetString("WEB_SERVER_PORT"))
		webServer.Start()
	}()

	select {
	case <-signCh:
		log.Println("Shutting down gracefully")
	case <-ctx.Done():
		log.Println("Shutting down gracefully")
	}

	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
}

func initProvider(serviceName, collectorUrl string) (func(ctx context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, collectorUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return traceProvider.Shutdown, nil
}
