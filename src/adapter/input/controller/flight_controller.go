package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/adapter/input/model/response"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/configuration/logger"
	"github.com/jinzhu/copier"
	"net/http"
	"time"

	flightRequest "github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/adapter/input/model/request"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/application/domain"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/application/port/input"

	"github.com/aws/aws-lambda-go/events"
)

var FlightControllerInput input.FlightUseCase

func FlightController(ctx context.Context, event *flightRequest.FlightRequest) (
	events.APIGatewayProxyResponse, error) {

	logger.Info(
		fmt.Sprintf("method=FlightController, action=init, request=%+v", event))

	if len(event.Departure) != 3 || len(event.Arrival) != 3 {
		logger.Error("method=FlightController, action=error", errors.New("invalid departure or arrival"))
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}, errors.New("invalid departure or arrival fields")
	}

	departureTime, err := time.Parse(time.RFC3339, event.DepartureDateTime)
	if err != nil {
		logger.Error("method=FlightController, action=error", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}, errors.New(fmt.Sprintf("error trying to parse departureTime datetime, error=%s", err.Error()))
	}

	arrivalTime, err := time.Parse(time.RFC3339, event.ArrivalDateTime)
	if err != nil {
		logger.Error("method=FlightController, action=error", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}, errors.New(fmt.Sprintf("error trying to parse arrivalTime datetime, error=%s", err.Error()))
	}

	flightDomain := domain.FlightDomain{
		Departure:         event.Departure,
		DepartureDateTime: departureTime,
		Arrival:           event.Arrival,
		ArrivalDateTime:   arrivalTime,
	}

	flightResponse, errRes := FlightControllerInput.GetConnectingFlights(ctx, flightDomain)
	if errRes != nil {
		logger.Error("method=FlightController, action=error", err)
		return events.APIGatewayProxyResponse{
			StatusCode: errRes.Code,
			Body:       errRes.Message,
		}, err
	}

	var lambdaFlightResponse []response.LambdaConnectingFlightsResponse

	err = copier.Copy(&lambdaFlightResponse, &flightResponse)
	if err != nil {
		logger.Error("method=FlightController, action=error", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, err
	}

	responseBody, errMarshal := json.Marshal(lambdaFlightResponse)
	if errMarshal != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error trying to parse the flightConnection response to json, error=%s", errMarshal),
		}, errMarshal
	}

	logger.Info(
		fmt.Sprintf("method=FlightController, action=success, response=%+v", flightResponse))
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBody),
	}, nil
}
