FROM golang:alpine

WORKDIR /app

COPY . .
COPY ./go.mod .

RUN go mod download

RUN go build -o app cmd/main.go

CMD ["./app"]