FROM golang:1.21.3 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go mod download

RUN go build -o /out/grafana-sync cmd/grafana-sync/main.go



FROM golang:1.21.3 AS release

WORKDIR /app

COPY --from=build /out/grafana-sync /app/

ENTRYPOINT ["/app/grafana-sync"]
