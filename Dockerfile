FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o stress-test

ENTRYPOINT ["./stress-test"]
