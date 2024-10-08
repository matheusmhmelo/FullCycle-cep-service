// Code generated by MockGen. DO NOT EDIT.
// Source: ./cep_api/internal/usecase/weather.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	usecase "github.com/matheusmhmelo/FullCycle-cep-api/internal/usecase"
)

// MockWeatherUseCase is a mock of WeatherUseCase interface.
type MockWeatherUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockWeatherUseCaseMockRecorder
}

// MockWeatherUseCaseMockRecorder is the mock recorder for MockWeatherUseCase.
type MockWeatherUseCaseMockRecorder struct {
	mock *MockWeatherUseCase
}

// NewMockWeatherUseCase creates a new mock instance.
func NewMockWeatherUseCase(ctrl *gomock.Controller) *MockWeatherUseCase {
	mock := &MockWeatherUseCase{ctrl: ctrl}
	mock.recorder = &MockWeatherUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWeatherUseCase) EXPECT() *MockWeatherUseCaseMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockWeatherUseCase) Execute(ctx context.Context, cep string) (*usecase.Weather, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx, cep)
	ret0, _ := ret[0].(*usecase.Weather)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockWeatherUseCaseMockRecorder) Execute(ctx, cep interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockWeatherUseCase)(nil).Execute), ctx, cep)
}
