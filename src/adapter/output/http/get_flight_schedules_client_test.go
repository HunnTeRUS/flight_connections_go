package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFlightSchedules_Success(t *testing.T) {
	client := NewFlightClient()

	scheduleResponse := `{"month": 6, "days": [{"day": 1, "flights": [{"number": "1926", "departureTime": "18:00", "arrivalTime": "21:35"}]}]}`
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(scheduleResponse))
	}))
	defer mockServer.Close()

	os.Setenv(RYANAIR_SCHEDULES_API, mockServer.URL+"/%s/%s/%d/%d")

	schedule, err := client.GetFlightSchedules("DUB", "WRO", 2023, 9)

	assert.Nil(t, err)
	assert.NotNil(t, schedule)
	assert.Equal(t, 6, schedule.Month)
	assert.Equal(t, 1, len(schedule.Days))
	assert.Equal(t, 1, len(schedule.Days[0].Flights))
	assert.Equal(t, "1926", schedule.Days[0].Flights[0].Number)
	assert.Equal(t, "18:00", schedule.Days[0].Flights[0].DepartureTime)
	assert.Equal(t, "21:35", schedule.Days[0].Flights[0].ArrivalTime)
}
