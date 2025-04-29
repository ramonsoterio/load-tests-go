FROM golang:1.23.2

WORKDIR /app

COPY go.mod ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/stress_test .

ENTRYPOINT ["/app/stress_test"]