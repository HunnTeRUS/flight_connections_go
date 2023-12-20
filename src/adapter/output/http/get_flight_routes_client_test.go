// http_test.go
package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllFlightRoutes_Success(t *testing.T) {
	client := NewFlightClient()

	routesResponse := `[{"airportFrom":"AAL","airportTo":"STN", "connectingAirport":null,"newRoute":false,"seasonalRoute":false,"operator":"RYANAIR","carrierCode":"FR","group":"CITY","similarArrivalAirportCodes":[],"tags":[]}]`
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(routesResponse))
	}))
	defer mockServer.Close()

	os.Setenv(RYANAIR_ROUTES_API, mockServer.URL)

	routes, err := client.GetAllFlightRoutes()

	assert.Nil(t, err)
	assert.NotNil(t, routes)
	assert.Equal(t, 1, len(routes))
	assert.Equal(t, "AAL", routes[0].AirportFrom)
	assert.Equal(t, "STN", routes[0].AirportTo)
}

func TestGetAllFlightRoutes_InvalidJSONResponse(t *testing.T) {
	client := NewFlightClient()

	invalidJSONResponse := `{"invalid": "json"}`
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(invalidJSONResponse))
	}))
	defer mockServer.Close()

	os.Setenv(RYANAIR_ROUTES_API, mockServer.URL)

	routes, err := client.GetAllFlightRoutes()

	assert.Nil(t, routes)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Code)
	assert.Contains(t, err.Message, "Error trying to read and process Ryanair routes API response")
}
