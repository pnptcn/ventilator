FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .
RUN go build -o ventilator

FROM bitnami/minideb:latest

WORKDIR /app

RUN install_packages ca-certificates \
	&& useradd -m -s /bin/bash nonroot \
	&& chown -R nonroot:nonroot /app

COPY --from=builder /app/ventilator .

USER nonroot

EXPOSE 8080

CMD ./ventilator serve
