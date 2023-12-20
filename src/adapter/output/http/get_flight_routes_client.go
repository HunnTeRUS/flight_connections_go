package http

import (
	"encoding/json"
	"fmt"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/adapter/output/model/response"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/application/domain"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/configuration/logger"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/configuration/rest_err"
	"github.com/jinzhu/copier"
	"io"
	"net/http"
	"os"
)

const (
	RYANAIR_SCHEDULES_API = "RYANAIR_SCHEDULES_API"
	RYANAIR_ROUTES_API    = "RYANAIR_ROUTES_API"
)

func NewFlightClient() *FlightClient {
	return &FlightClient{}
}

type FlightClient struct {
}

func (*FlightClient) GetAllFlightRoutes() ([]domain.RoutesDomain, *rest_err.RestErr) {
	logger.Info("method=GetAllFlightRoutes, action=init")

	resp, err := http.Get(os.Getenv(RYANAIR_ROUTES_API))
	if err != nil {
		logger.Error("method=GetAllFlightRoutes, action=error", err)
		return nil, rest_err.NewInternalServerError(
			fmt.Sprintf("Error trying to call Ryanair routes API, error=%s", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("method=GetAllFlightRoutes, action=error", err)
		return nil, rest_err.NewInternalServerError(
			fmt.Sprintf("Error trying to read and process Ryanair routes API response, error=%s", err))
	}

	var routes []response.Route
	if err := json.Unmarshal(body, &routes); err != nil {
		logger.Error("method=GetAllFlightRoutes, action=error", err)
		return nil, rest_err.NewInternalServerError(
			fmt.Sprintf("Error trying to read and process Ryanair routes API response, error=%s", err))
	}

	var routeDomain []domain.RoutesDomain
	if err = copier.Copy(&routeDomain, &routes); err != nil {
		logger.Error("method=GetAllFlightRoutes, action=error", err)
		return nil, rest_err.NewInternalServerError(
			fmt.Sprintf("Error trying copy data, error=%s", err))
	}

	logger.Info(fmt.Sprintf("method=GetAllFlightRoutes, action=success, numOfFlights=%d", len(routeDomain)))
	return routeDomain, nil
}
