package domain

import "time"

type FlightDomain struct {
	Departure         string
	DepartureDateTime time.Time
	Arrival           string
	ArrivalDateTime   time.Time
}

type RoutesDomain struct {
	AirportFrom       string `copier:"AirportFrom"`
	AirportTo         string `copier:"AirportTo"`
	ConnectingAirport string `copier:"ConnectingAirport"`
	Operator          string `copier:"Operator"`
}

type FlightScheduleDomain struct {
	Number        string `copier:"Number"`
	DepartureTime string `copier:"DepartureTime"`
	ArrivalTime   string `copier:"ArrivalTime"`
}

type ScheduleDomain struct {
	Month int `copier:"Month"`
	Days  []struct {
		Day     int                    `copier:"Day"`
		Flights []FlightScheduleDomain `copier:"Flights"`
	} `copier:"Days"`
}

type ConnectingFlightsDomain struct {
	Stops int `copier:"stops"`
	Legs  []struct {
		DepartureAirport  string `copier:"DepartureAirport"`
		ArrivalAirport    string `copier:"ArrivalAirport"`
		DepartureDateTime string `copier:"DepartureDateTime"`
		ArrivalDateTime   string `copier:"ArrivalDateTime"`
	} `copier:"legs"`
}
