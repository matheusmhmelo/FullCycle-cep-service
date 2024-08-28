package gateway

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

func TestWeatherGatewayImpl_ValidateLocation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		getFunc = func(_ string) (resp *http.Response, err error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"localidade": "City"}`)),
			}, nil
		}
		defer func() {
			getFunc = http.Get
		}()

		gt := &weatherGatewayImpl{}
		got, err := gt.ValidateLocation("99999999")
		require.NoError(t, err)
		require.Equal(t, "City", gt.location)
		require.Equal(t, "City", got)
	})
	t.Run("empty cep", func(t *testing.T) {
		gt := &weatherGatewayImpl{}
		_, err := gt.ValidateLocation("")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrorInvalidCEP)
	})
	t.Run("response error", func(t *testing.T) {
		getFunc = func(_ string) (resp *http.Response, err error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"erro": "error"}`)),
			}, nil
		}
		defer func() {
			getFunc = http.Get
		}()

		gt := &weatherGatewayImpl{}
		_, err := gt.ValidateLocation("99999999")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrorNotFoundCEP)
	})
}

func TestWeatherGatewayImpl_GetWeather(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		getFunc = func(_ string) (resp *http.Response, err error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"current": {"temp_c": 1}}`)),
			}, nil
		}
		defer func() {
			getFunc = http.Get
		}()

		gt := &weatherGatewayImpl{}
		got, err := gt.GetWeather()
		require.NoError(t, err)
		require.Equal(t, float64(1), got)
	})
	t.Run("response error", func(t *testing.T) {
		getFunc = func(_ string) (resp *http.Response, err error) {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
			}, nil
		}
		defer func() {
			getFunc = http.Get
		}()

		gt := &weatherGatewayImpl{}
		_, err := gt.GetWeather()
		require.Error(t, err)
	})
}
