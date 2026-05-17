FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service ./cmd/service

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -a github.com/pressly/goose/v3/cmd/goose@latest
FROM alpine:3.23



COPY --from=builder /app/service ./
COPY --from=builder /app/config/ ./config/
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/migrations/ ./migrations/
COPY --from=builder /app/scripts/entrypoint.sh ./

RUN chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]