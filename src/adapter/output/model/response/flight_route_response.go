package response

type Route struct {
	AirportFrom       string `json:"airportFrom" copier:"AirportFrom"`
	AirportTo         string `json:"airportTo" copier:"AirportTo"`
	ConnectingAirport string `json:"connectingAirport" copier:"ConnectingAirport"`
	Operator          string `json:"operator" copier:"Operator"`
}
