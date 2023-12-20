FROM golang:1.19 AS BUILDER

WORKDIR /app
COPY src src
COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o ryanairflights .

FROM public.ecr.aws/lambda/go:1 as RUNNER

ENV _LAMBDA_SERVER_PORT 9000
ENV AWS_LAMBDA_RUNTIME_API localhost:8080

ENV RYANAIR_ROUTES_API https://services-api.ryanair.com/views/locate/3/routes
ENV RYANAIR_SCHEDULES_API https://services-api.ryanair.com/timtbl/3/schedules/%s/%s/years/%d/months/%d
ENV MAX_INTERCONNECT_WAIT 2h
ENV TIME_LAYOUT 2006-01-02T15:04

COPY --from=BUILDER /app/ryanairflights ${LAMBDA_TASK_ROOT}
COPY --from=BUILDER /app/ryanairflights /app/ryanairflights

ENTRYPOINT ["/app/ryanairflights"]