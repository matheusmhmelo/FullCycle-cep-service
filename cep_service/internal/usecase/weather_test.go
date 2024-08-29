package usecase

import (
	"bytes"
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

func TestWeatherUseCaseImpl_Execute(t *testing.T) {
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

		gt := &weatherUseCaseImpl{
			baseURL: "http://localhost.com:8080",
		}
		got, code, err := gt.Execute(context.Background(), "99999999")
		require.NoError(t, err)
		require.Equal(t, []byte(`{"localidade": "City"}`), got)
		require.Equal(t, code, http.StatusOK)
	})
	t.Run("response with error", func(t *testing.T) {
		doFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Body:       io.NopCloser(bytes.NewBufferString("invalid CEP received")),
			}, nil
		}
		defer func() {
			doFunc = http.DefaultClient.Do
		}()

		gt := &weatherUseCaseImpl{
			baseURL: "http://localhost.com:8080",
		}
		got, code, err := gt.Execute(context.Background(), "99999999")
		require.NoError(t, err)
		require.Equal(t, []byte("invalid CEP received"), got)
		require.Equal(t, code, http.StatusUnprocessableEntity)
	})
	t.Run("error to do request", func(t *testing.T) {
		doFunc = func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("error")
		}
		defer func() {
			doFunc = http.DefaultClient.Do
		}()

		gt := &weatherUseCaseImpl{
			baseURL: "",
		}
		_, code, err := gt.Execute(context.Background(), "99999999")
		require.Error(t, err)
		require.Equal(t, code, http.StatusInternalServerError)
	})
}
