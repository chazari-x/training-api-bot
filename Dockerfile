FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:3.10

RUN adduser -DH training-api

WORKDIR /app

COPY --from=builder /app/main /app/

COPY etc/config.yaml /app/etc/config.yaml
RUN chown training-api:training-api /app
RUN chmod +x /app

USER training-api

CMD ["/app/main"]