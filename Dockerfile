FROM golang:1.21
LABEL authors="guvanch"
WORKDIR /app
COPY . .
EXPOSE 8085
RUN go mod tidy
RUN cd cmd/
RUN go build -o bin/app
RUN cd bin
CMD ["./app"]
