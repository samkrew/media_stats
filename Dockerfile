FROM golang:1.8-alpine
MAINTAINER Semyon Krivosheev

ENV FFMPEG_VERSION=3.0.2

RUN apk add --update ffmpeg-dev libc-dev gcc pkgconfig && rm -rf /var/cache/apk/*

ADD . /go/src/github.com/samkrew/media_stats
RUN go install -tags ffmpeg33 github.com/samkrew/media_stats

ENTRYPOINT ["/go/bin/media_stats"]
#ENTRYPOINT ["/bin/sh"]

# Document that the service listens on port 8080.
EXPOSE 3000