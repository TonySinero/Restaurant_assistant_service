FROM golang:1.17-alpine3.15 AS builder

COPY . /restaurant-assistant/
WORKDIR /restaurant-assistant/

RUN go mod download
RUN go build -o ./bin/app cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /restaurant-assistant/bin/app .
COPY --from=builder /restaurant-assistant/configs configs/
COPY --from=builder /restaurant-assistant/migrations migrations/

EXPOSE 80 50080

CMD ["./app"]