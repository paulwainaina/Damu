FROM golang:latest

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy

RUN go build -o ./app

EXPOSE 8080

CMD ["./app"]