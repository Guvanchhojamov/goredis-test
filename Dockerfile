FROM golang:1.21
LABEL authors="guvanch"
WORKDIR /app
COPY . .
RUN go mod tidy