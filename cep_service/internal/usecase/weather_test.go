package usecase

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/matheusmhmelo/FullCycle-cep-service/cep_api/internal/infra/gateway/mock_gateway"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWeatherUseCaseImpl_Execute(t *testing.T) {
	tests := []struct {
		name         string
		prepareMocks func(t *testing.T) *mock_gateway.MockWeatherGatewayInterface
		assertFunc   func(t *testing.T, got *Weather, err error)
	}{
		{
			name: "success",
			prepareMocks: func(t *testing.T) *mock_gateway.MockWeatherGatewayInterface {
				ctrl := gomock.NewController(t)
				mock := mock_gateway.NewMockWeatherGatewayInterface(ctrl)
				mock.EXPECT().ValidateLocation("cep").Return("city", nil).Times(1)
				mock.EXPECT().GetWeather().Return(float64(1), nil).Times(1)
				return mock
			},
			assertFunc: func(t *testing.T, got *Weather, err error) {
				require.NoError(t, err)
				require.Equal(t, &Weather{
					City:       "city",
					Fahrenheit: 33.8,
					Celsius:    1,
					Kelvin:     274,
				}, got)
			},
		},
		{
			name: "error validating location",
			prepareMocks: func(t *testing.T) *mock_gateway.MockWeatherGatewayInterface {
				ctrl := gomock.NewController(t)
				mock := mock_gateway.NewMockWeatherGatewayInterface(ctrl)
				mock.EXPECT().ValidateLocation("cep").Return("", errors.New("error")).Times(1)
				return mock
			},
			assertFunc: func(t *testing.T, got *Weather, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "error to get weather",
			prepareMocks: func(t *testing.T) *mock_gateway.MockWeatherGatewayInterface {
				ctrl := gomock.NewController(t)
				mock := mock_gateway.NewMockWeatherGatewayInterface(ctrl)
				mock.EXPECT().ValidateLocation("cep").Return("", nil).Times(1)
				mock.EXPECT().GetWeather().Return(float64(0), errors.New("error")).Times(1)
				return mock
			},
			assertFunc: func(t *testing.T, got *Weather, err error) {
				require.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.prepareMocks(t)
			usecase := NewWeatherUseCase(mock)
			got, err := usecase.Execute("cep")
			tt.assertFunc(t, got, err)
		})
	}
}
