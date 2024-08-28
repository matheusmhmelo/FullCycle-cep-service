package web

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/matheusmhmelo/FullCycle-cep-service/cep_api/internal/usecase"
	"github.com/matheusmhmelo/FullCycle-cep-service/cep_api/internal/usecase/mock_usecase"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOrderHandler_Get(t *testing.T) {
	expected := &usecase.Weather{
		Fahrenheit: 10,
		Celsius:    5,
		Kelvin:     1,
	}

	ctrl := gomock.NewController(t)
	mock := mock_usecase.NewMockWeatherUseCase(ctrl)
	mock.EXPECT().Execute("cep").Return(expected, nil).Times(1)
	handler := OrderHandler{
		weather: mock,
	}

	req := httptest.NewRequest(http.MethodGet, "/test?cep=cep", nil)
	w := httptest.NewRecorder()
	handler.Get(w, req)

	res := w.Result()
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)

	var got usecase.Weather
	err := json.NewDecoder(res.Body).Decode(&got)
	require.NoError(t, err)
	require.Equal(t, expected, &got)
}
