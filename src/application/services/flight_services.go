package services

import (
	"context"
	"fmt"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/application/domain"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/application/port/output"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/configuration/logger"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/configuration/rest_err"
	"os"
	"time"
)

var (
	TIME_LAYOUT           = ""
	MAX_INTERCONNECT_WAIT time.Duration
)

func NewFlightServices(flightPort output.FlightPort) *flightServices {
	TIME_LAYOUT = os.Getenv("TIME_LAYOUT")
	MAX_INTERCONNECT_WAIT = ParseDuration(os.Getenv("MAX_INTERCONNECT_WAIT"))

	return &flightServices{flightPort}
}

type flightServices struct {
	flightPort output.FlightPort
}

func (fs *flightServices) GetConnectingFlights(ctx context.Context, flightDomain domain.FlightDomain) (
	[]domain.ConnectingFlightsDomain, *rest_err.RestErr) {

	logger.Info(
		fmt.Sprintf("method=GetConnectingFlights, action=init, domain=%+v", flightDomain))

	// TODO: Use the context to apply some timeouts or traceability
	routes, err := fs.flightPort.GetAllFlightRoutes()
	if err != nil {
		return nil, err
	}

	var results []domain.ConnectingFlightsDomain

	// Direct
	for _, route := range routes {
		if route.AirportFrom == flightDomain.Departure && route.AirportTo == flightDomain.Arrival {
			schedule, err := fs.flightPort.GetFlightSchedules(
				flightDomain.Departure,
				flightDomain.Arrival,
				flightDomain.DepartureDateTime.Year(),
				int(flightDomain.DepartureDateTime.Month()))
			if err != nil || schedule.Days == nil {
				logger.Error("method=GetConnectingFlights, action=error", err)
				return nil, err
			}

			for _, day := range schedule.Days {
				if day.Day >= flightDomain.DepartureDateTime.Day() {
					for _, flight := range day.Flights {
						departureTime, _ := time.Parse(
							TIME_LAYOUT,
							fmt.Sprintf("%d-%02d-%02dT%s",
								flightDomain.DepartureDateTime.Year(),
								flightDomain.DepartureDateTime.Month(),
								day.Day,
								flight.DepartureTime))
						arrivalTime, _ := time.Parse(TIME_LAYOUT,
							fmt.Sprintf("%d-%02d-%02dT%s",
								flightDomain.DepartureDateTime.Year(),
								flightDomain.DepartureDateTime.Month(),
								day.Day,
								flight.ArrivalTime))

						if departureTime.After(flightDomain.DepartureDateTime) && arrivalTime.Before(flightDomain.ArrivalDateTime) {
							results = append(results, domain.ConnectingFlightsDomain{
								Stops: 0,
								Legs: []struct {
									DepartureAirport  string `copier:"DepartureAirport"`
									ArrivalAirport    string `copier:"ArrivalAirport"`
									DepartureDateTime string `copier:"DepartureDateTime"`
									ArrivalDateTime   string `copier:"ArrivalDateTime"`
								}{
									{
										DepartureAirport:  flightDomain.Departure,
										ArrivalAirport:    flightDomain.Arrival,
										DepartureDateTime: departureTime.Format(TIME_LAYOUT),
										ArrivalDateTime:   arrivalTime.Format(TIME_LAYOUT),
									},
								},
							})
						}
					}
				}
			}
		} else if route.AirportFrom == flightDomain.Departure && route.ConnectingAirport != "" && route.AirportTo == route.ConnectingAirport {
			// Interconnected flight
			firstLegSchedule, err := fs.flightPort.GetFlightSchedules(
				flightDomain.Departure,
				flightDomain.Arrival,
				flightDomain.DepartureDateTime.Year(),
				int(flightDomain.DepartureDateTime.Month()))
			if err != nil {
				logger.Error("method=GetConnectingFlights, action=error", err)
				return nil, err
			}

			for _, firstLegDay := range firstLegSchedule.Days {

				// TODO: Transform each schedule call to the GetFlightSchedules in a goroutine, so we can increase performance
				if firstLegDay.Day >= flightDomain.DepartureDateTime.Day() {
					for _, firstLegFlight := range firstLegDay.Flights {
						firstLegDepartureTime, _ := time.Parse(
							TIME_LAYOUT,
							fmt.Sprintf("%d-%02d-%02dT%s",
								flightDomain.DepartureDateTime.Year(),
								flightDomain.DepartureDateTime.Month(),
								firstLegDay.Day,
								firstLegFlight.DepartureTime))
						firstLegArrivalTime, _ := time.Parse(
							TIME_LAYOUT,
							fmt.Sprintf("%d-%02d-%02dT%s",
								flightDomain.DepartureDateTime.Year(),
								flightDomain.DepartureDateTime.Month(),
								firstLegDay.Day,
								firstLegFlight.ArrivalTime))

						if firstLegDepartureTime.After(flightDomain.DepartureDateTime) {
							secondLegSchedule, err := fs.flightPort.GetFlightSchedules(
								route.ConnectingAirport,
								flightDomain.Arrival,
								firstLegArrivalTime.Year(),
								int(firstLegArrivalTime.Month()))
							if err != nil {
								logger.Error("method=GetConnectingFlights, action=error", err)
								return nil, err
							}

							for _, secondLegDay := range secondLegSchedule.Days {
								if secondLegDay.Day >= firstLegArrivalTime.Day() {
									for _, secondLegFlight := range secondLegDay.Flights {
										secondLegDepartureTime, _ := time.Parse(
											TIME_LAYOUT,
											fmt.Sprintf("%d-%02d-%02dT%s",
												firstLegArrivalTime.Year(),
												firstLegArrivalTime.Month(),
												secondLegDay.Day,
												secondLegFlight.DepartureTime))
										secondLegArrivalTime, _ := time.Parse(
											TIME_LAYOUT,
											fmt.Sprintf("%d-%02d-%02dT%s",
												firstLegArrivalTime.Year(),
												firstLegArrivalTime.Month(),
												secondLegDay.Day,
												secondLegFlight.ArrivalTime))

										if secondLegDepartureTime.Sub(firstLegArrivalTime) >= MAX_INTERCONNECT_WAIT &&
											secondLegArrivalTime.Before(flightDomain.ArrivalDateTime) {
											results = append(results, domain.ConnectingFlightsDomain{
												Stops: 1,
												Legs: []struct {
													DepartureAirport  string `copier:"DepartureAirport"`
													ArrivalAirport    string `copier:"ArrivalAirport"`
													DepartureDateTime string `copier:"DepartureDateTime"`
													ArrivalDateTime   string `copier:"ArrivalDateTime"`
												}{
													{
														DepartureAirport:  flightDomain.Departure,
														ArrivalAirport:    route.ConnectingAirport,
														DepartureDateTime: firstLegDepartureTime.Format(TIME_LAYOUT),
														ArrivalDateTime:   firstLegArrivalTime.Format(TIME_LAYOUT),
													},
													{
														DepartureAirport:  route.ConnectingAirport,
														ArrivalAirport:    flightDomain.Arrival,
														DepartureDateTime: secondLegDepartureTime.Format(TIME_LAYOUT),
														ArrivalDateTime:   secondLegArrivalTime.Format(TIME_LAYOUT),
													},
												},
											})
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	logger.Info(
		fmt.Sprintf("method=GetConnectingFlights, action=success, domain=%+v", results))

	return results, nil
}

func ParseDuration(timeValue string) time.Duration {
	duration, _ := time.ParseDuration(timeValue)
	return duration
}
