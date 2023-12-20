package response

type Flight struct {
	Number        string `json:"number" copier:"Number"`
	DepartureTime string `json:"departureTime" copier:"DepartureTime"`
	ArrivalTime   string `json:"arrivalTime" copier:"ArrivalTime"`
}

type Schedule struct {
	Month int `json:"month" copier:"Month"`
	Days  []struct {
		Day     int      `json:"day" copier:"Day"`
		Flights []Flight `json:"flights" copier:"Flights"`
	} `json:"days" copier:"Days"`
}
