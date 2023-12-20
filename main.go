package main

import (
	"context"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/adapter/input/controller"
	flightRequest "github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/adapter/input/model/request"
	output "github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/adapter/output/http"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/application/services"
	"github.com/RyanairLabs/Java-Challenge---Otavio-Santos/src/configuration/logger"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

func main() {
	logger.Info("About to start application")
	godotenv.Load()

	flightClient := output.NewFlightClient()
	flightService := services.NewFlightServices(flightClient)
	controller.FlightControllerInput = flightService

	lambda.Start(FlightController)
}

func FlightController(ctx context.Context, event *flightRequest.FlightRequest) (
	events.APIGatewayProxyResponse, error) {

	return controller.FlightController(ctx, event)
}
