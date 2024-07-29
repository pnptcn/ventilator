FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go install github.com/air-verse/air@latest
COPY . .

EXPOSE 8080

CMD /go/bin/air
