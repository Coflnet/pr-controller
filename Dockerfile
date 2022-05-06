FROM golang:1.18.1-bullseye as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build cmd/pr-controller/main.go

CMD /app/main