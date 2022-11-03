FROM golang:1.17-alpine AS builder

WORKDIR /app

COPY . .

# install timzone data
RUN apk add --no-cache tzdata

RUN go mod tidy -compat="1.17"
RUN go mod download
RUN go build -v -o /go/bin/valbot

FROM ubuntu:latest AS installer

RUN apt-get -y update
RUN apt-get -y upgrade
RUN apt-get -y install ffmpeg
RUN ffmpeg -version
RUN which ffmpeg

FROM alpine:latest

COPY --from=installer /usr/bin/ffmpeg /usr/bin/ffmpeg
COPY --from=installer /usr/bin/ffmpeg /usr/local/bin/ffmpeg
RUN export PATH=/usr/local/bin:$PATH
RUN ffmepg -version

# copy timezone info from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /go/bin/valbot /valbot
RUN chmod +x /valbot

EXPOSE 8080

ENV TZ=Asia/Hong_Kong

WORKDIR /app

VOLUME [ "/app/data" ]

ENTRYPOINT [ "/valbot" ]