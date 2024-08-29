package web

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/matheusmhmelo/FullCycle-cep-service/internal/usecase/mock_usecase"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOrderHandler_Post(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := []byte(`{"test": "response"}`)

		ctrl := gomock.NewController(t)
		mock := mock_usecase.NewMockWeatherUseCase(ctrl)
		mock.EXPECT().Execute(gomock.Any(), "99999999").Return(expected, http.StatusOK, nil).Times(1)
		handler := OrderHandler{
			weather:    mock,
			otelTracer: otel.Tracer("cep-service"),
		}

		bytesObj := []byte(`{"cep": "99999999"}`)
		body := bytes.NewBuffer(bytesObj)

		req := httptest.NewRequest(http.MethodPost, "/test", body)
		w := httptest.NewRecorder()
		handler.Post(w, req)

		res := w.Result()
		defer res.Body.Close()
		require.Equal(t, http.StatusOK, res.StatusCode)

		got, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, expected, got)
	})
	t.Run("invalid zipcode", func(t *testing.T) {
		handler := OrderHandler{
			otelTracer: otel.Tracer("cep-service"),
		}

		bytesObj := []byte(`{"cep": "999999999999"}`)
		body := bytes.NewBuffer(bytesObj)

		req := httptest.NewRequest(http.MethodPost, "/test", body)
		w := httptest.NewRecorder()
		handler.Post(w, req)

		res := w.Result()
		defer res.Body.Close()
		require.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)

		got, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, "invalid zipcode\n", string(got))
	})
}
