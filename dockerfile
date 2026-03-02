FROM golang:1.24.2-alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN go build -o moneyServic cmd/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /build/moneyServic /app/moneyServic
COPY --from=builder /build/config.yml /app/config.yml
COPY --from=builder /build/migrations /app/migrations
COPY --from=builder  /build/public.pem /app/public.pem

CMD ["./moneyServic"]