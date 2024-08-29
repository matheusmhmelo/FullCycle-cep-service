package gateway

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"io"
	"net/http"
	"testing"
)

func TestWeatherGatewayImpl_ValidateLocation(t *testing.T) {
	traceMock := otel.Tracer("test")

	t.Run("success", func(t *testing.T) {
		doFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"localidade": "City"}`)),
			}, nil
		}
		defer func() {
			doFunc = http.DefaultClient.Do
		}()

		gt := &weatherGatewayImpl{
			otelTracer: traceMock,
		}
		got, err := gt.ValidateLocation(context.Background(), "99999999")
		require.NoError(t, err)
		require.Equal(t, "City", gt.location)
		require.Equal(t, "City", got)
	})
	t.Run("empty cep", func(t *testing.T) {
		gt := &weatherGatewayImpl{
			otelTracer: traceMock,
		}
		_, err := gt.ValidateLocation(context.Background(), "")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrorInvalidCEP)
	})
	t.Run("response error", func(t *testing.T) {
		doFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"erro": "error"}`)),
			}, nil
		}
		defer func() {
			doFunc = http.DefaultClient.Do
		}()

		gt := &weatherGatewayImpl{
			otelTracer: traceMock,
		}
		_, err := gt.ValidateLocation(context.Background(), "99999999")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrorNotFoundCEP)
	})
}

func TestWeatherGatewayImpl_GetWeather(t *testing.T) {
	traceMock := otel.Tracer("test")

	t.Run("success", func(t *testing.T) {
		doFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"current": {"temp_c": 1}}`)),
			}, nil
		}
		defer func() {
			doFunc = http.DefaultClient.Do
		}()

		gt := &weatherGatewayImpl{
			otelTracer: traceMock,
		}
		got, err := gt.GetWeather(context.Background())
		require.NoError(t, err)
		require.Equal(t, float64(1), got)
	})
	t.Run("response error", func(t *testing.T) {
		doFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
			}, nil
		}
		defer func() {
			doFunc = http.DefaultClient.Do
		}()

		gt := &weatherGatewayImpl{
			otelTracer: traceMock,
		}
		_, err := gt.GetWeather(context.Background())
		require.Error(t, err)
	})
}
