FROM golang:1.8-alpine
MAINTAINER Semyon Krivosheev

ENV FFMPEG_VERSION=3.0.2

RUN apk add --update ffmpeg-dev libc-dev gcc pkgconfig git && rm -rf /var/cache/apk/*

RUN go get -u github.com/kardianos/govendor

ADD . /go/src/github.com/samkrew/media_stats
RUN cd /go/src/github.com/samkrew/media_stats && govendor sync
RUN go install -tags ffmpeg33 github.com/samkrew/media_stats

ENTRYPOINT ["/go/bin/media_stats"]

EXPOSE 3000