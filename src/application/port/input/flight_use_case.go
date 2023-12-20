package input

import (
	"context"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/application/domain"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/configuration/rest_err"
)

type FlightUseCase interface {
	GetConnectingFlights(ctx context.Context, flightDomain domain.FlightDomain) (
		[]domain.ConnectingFlightsDomain, *rest_err.RestErr)
}
