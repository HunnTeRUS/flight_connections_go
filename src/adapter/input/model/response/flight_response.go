package response

type LambdaConnectingFlightsResponse struct {
	Stops int ` json:"stops" copier:"Stops"`
	Legs  []struct {
		DepartureAirport  string `json:"departureAirport" copier:"DepartureAirport"`
		ArrivalAirport    string `json:"arrivalAirport" copier:"ArrivalAirport"`
		DepartureDateTime string `json:"departureDateTime" copier:"DepartureDateTime"`
		ArrivalDateTime   string `json:"arrivalDateTime" copier:"ArrivalDateTime"`
	} `json:"legs" copier:"Legs"`
}
