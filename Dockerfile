FROM golang:1.17-alpine AS builder

WORKDIR /app

COPY . .

# install timzone data
RUN apk add --no-cache tzdata

RUN go mod tidy -compat="1.17"
RUN go mod download
RUN go build -v -o /go/bin/valbot

FROM alpine:latest

# copy timezone info from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /go/bin/valbot /valbot
RUN chmod +x /valbot

EXPOSE 8080

ENV TZ=Asia/Hong_Kong

WORKDIR /app

VOLUME [ "/app/data" ]

ENTRYPOINT [ "/valbot" ]