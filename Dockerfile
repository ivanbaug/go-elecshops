# syntax=docker/dockerfile:1

## Build
FROM golang:1.20 AS build
#FROM golang:1.20-bullseye

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY ./cmd ./cmd
COPY ./internal ./internal

RUN go mod download

RUN CGO_ENABLED=0 go build -o /go-elecsh cmd/main.go

## Deploy
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /go-elecsh /go-elecsh
COPY .env ./

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/go-elecsh"]
