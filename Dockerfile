FROM docker.io/golang:1.23-alpine3.21 AS builder

WORKDIR /code

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY internal ./internal
RUN go build -o /ssh_key_exporter

FROM scratch

COPY --from=builder /ssh_key_exporter /ssh_key_exporter
CMD [ "/ssh_key_exporter" ]
