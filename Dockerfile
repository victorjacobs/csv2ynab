FROM golang:1.21-alpine AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/csv2ynab

FROM scratch
COPY --from=builder /out/csv2ynab /csv2ynab
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/csv2ynab", "-watch", "-config", "./ynab.json"]
