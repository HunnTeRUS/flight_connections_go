package http

import (
	"encoding/json"
	"fmt"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/adapter/output/model/response"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/application/domain"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/configuration/rest_err"
	"github.com/jinzhu/copier"
	"io"
	"net/http"
	"os"
)

type localFlightSave struct {
	departure, arrival string
	year, month        int
}

var (
	searchedFlights = make(map[localFlightSave]*domain.ScheduleDomain)
)

func (*FlightClient) GetFlightSchedules(
	departure,
	arrival string,
	year,
	month int) (*domain.ScheduleDomain, *rest_err.RestErr) {
	url := fmt.Sprintf(
		os.Getenv(RYANAIR_SCHEDULES_API),
		departure,
		arrival,
		year,
		month)

	flight, ok := searchedFlights[localFlightSave{departure, arrival, year, month}]
	if ok {
		return flight, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(url)
		return nil, rest_err.NewInternalServerError(
			fmt.Sprintf("Error trying to call Ryanair routes API, error=%s", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(url)
		return nil, rest_err.NewInternalServerError(
			fmt.Sprintf("Error trying to read and process Ryanair routes API response, error=%s", err))
	}

	var schedule response.Schedule
	if err := json.Unmarshal(body, &schedule); err != nil {
		return nil, rest_err.NewInternalServerError(
			fmt.Sprintf("Error trying to read and process Ryanair routes API response, error=%s", err))
	}

	var scheduleDomain domain.ScheduleDomain
	if err := copier.Copy(&scheduleDomain, &schedule); err != nil {
		return nil, rest_err.NewInternalServerError(
			fmt.Sprintf("Error trying copy data, error=%s", err))
	}

	searchedFlights[localFlightSave{departure, arrival, year, month}] = &scheduleDomain

	return &scheduleDomain, nil
}
