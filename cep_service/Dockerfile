FROM golang:1.21

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/weather_service cmd/service/main.go

EXPOSE 8080

ENTRYPOINT ["/app/weather_service"]