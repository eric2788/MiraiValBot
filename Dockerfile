FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY . .

# install timzone data
RUN apk add --no-cache tzdata

RUN go mod tidy -compat="1.20"
RUN go mod download
RUN go build -v -o /go/bin/valbot

FROM alpine:latest

# no use for now
# RUN apk add --no-cache ffmpeg

# copy timezone info from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /go/bin/valbot /valbot
RUN chmod +x /valbot

EXPOSE 8080
EXPOSE 45678

ENV TZ=Asia/Hong_Kong
ENV COMPRESS_TYPE=zlib
ENV CACHE_STRATEGY=local

WORKDIR /app

VOLUME [ "/app/data", "/app/cache" ]

ENTRYPOINT [ "/valbot" ]