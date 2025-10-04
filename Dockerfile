FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download -x

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o notificationproxy 

RUN useradd -u 10001 appuser

FROM scratch

WORKDIR /app
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER appuser
COPY --chown=appuser:appuser --from=builder /app/notificationproxy /app/

EXPOSE 8080
EXPOSE 2525

ENTRYPOINT [ "/app/notificationproxy" ]