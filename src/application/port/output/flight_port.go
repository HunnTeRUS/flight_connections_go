package output

import (
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/application/domain"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/configuration/rest_err"
)

type FlightPort interface {
	GetFlightSchedules(departure, arrival string, year, month int) (
		*domain.ScheduleDomain, *rest_err.RestErr)
	GetAllFlightRoutes() ([]domain.RoutesDomain, *rest_err.RestErr)
}
