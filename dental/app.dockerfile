FROM golang:1.19.3-alpine3.16

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o mainApp ./cmd/api

# RUN go run ./cmd/api/
CMD [ "./mainApp" ]
