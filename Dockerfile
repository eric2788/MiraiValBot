FROM golang:1.17-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /go/bin/valbot

FROM alpine:latest

COPY --from=builder /go/bin/valbot /valbot
RUN chmod +x /valbot

EXPOSE 8080

WORKDIR /app

VOLUME [ "/app/data" ]

ENTRYPOINT [ "/valbot" ]