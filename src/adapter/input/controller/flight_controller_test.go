package controller

import (
	"context"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/adapter/input/model/request"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/application/domain"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/configuration/rest_err"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func TestFlightController_InvalidBody(t *testing.T) {
	event := createTestEvent("", "", "", "")

	response, err := FlightController(context.Background(), event)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestFlightController_InvalidJSON(t *testing.T) {
	event := createTestEvent("", "", "", "")

	response, err := FlightController(context.Background(), event)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestFlightController_InvalidDepartureTime(t *testing.T) {
	event := createTestEvent("ABD", "BCD", "invalid_time", "2023-01-01T14:00:00Z")

	response, err := FlightController(context.Background(), event)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, err.Error(), "error trying to parse departureTime datetime")
}

func TestFlightController_InvalidDepartureOrArrival(t *testing.T) {
	event := createTestEvent("", "B", "invalid_time", "2023-01-01T14:00:00Z")

	response, err := FlightController(context.Background(), event)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, err.Error(), "invalid departure or arrival fields")
}

func TestFlightController_InvalidArrivalTime(t *testing.T) {
	event := createTestEvent("ABC", "BCD", "2023-01-01T12:00:00Z", "invalid_time")

	response, err := FlightController(context.Background(), event)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, err.Error(), "error trying to parse arrivalTime datetime")
}

func TestFlightController_GetConnectingFlights_Failure(t *testing.T) {
	event := createTestEvent("ABC", "BCD", "2023-01-01T12:00:00Z", "2023-01-01T14:00:00Z")

	mockUseCase := new(MockUseCase)
	mockUseCase.On("GetConnectingFlights", mock.Anything, mock.Anything).Return(
		[]domain.ConnectingFlightsDomain{}, rest_err.NewInternalServerError("useCase error"))
	FlightControllerInput = mockUseCase

	response, err := FlightController(context.Background(), event)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.Contains(t, response.Body, "useCase error")
}

func TestFlightController_GetConnectingFlights_Success(t *testing.T) {
	event := createTestEvent("ABC", "BCD", "2023-01-01T12:00:00Z", "2023-01-01T14:00:00Z")

	mockUseCase := new(MockUseCase)
	mockUseCase.On("GetConnectingFlights", mock.Anything, mock.Anything).Return([]domain.ConnectingFlightsDomain{
		{
			Stops: 0,
			Legs: []struct {
				DepartureAirport  string `copier:"DepartureAirport"`
				ArrivalAirport    string `copier:"ArrivalAirport"`
				DepartureDateTime string `copier:"DepartureDateTime"`
				ArrivalDateTime   string `copier:"ArrivalDateTime"`
			}{
				{"ABD", "BCD", "2023-01-01T12:00:00Z", "2023-01-01T14:00:00Z"},
			},
		},
	}, nil)
	FlightControllerInput = mockUseCase

	response, err := FlightController(context.Background(), event)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

type MockUseCase struct {
	mock.Mock
}

func (m *MockUseCase) GetConnectingFlights(ctx context.Context, flightDomain domain.FlightDomain) (
	[]domain.ConnectingFlightsDomain, *rest_err.RestErr) {
	args := m.Called(ctx, flightDomain)
	arg1 := args.Get(1)
	if arg1 == nil {
		return args.Get(0).([]domain.ConnectingFlightsDomain), nil
	}

	return args.Get(0).([]domain.ConnectingFlightsDomain), args.Get(1).(*rest_err.RestErr)
}

func createTestEvent(departure, arrival, departureDateTime, arrivalDateTime string) *request.FlightRequest {
	return &request.FlightRequest{
		Departure:         departure,
		DepartureDateTime: departureDateTime,
		Arrival:           arrival,
		ArrivalDateTime:   arrivalDateTime,
	}
}
