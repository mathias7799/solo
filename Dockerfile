FROM golang:1.14-alpine as builder

ADD . /src/app

RUN cd /src/app && go build

FROM alpine:3

COPY --from=builder /src/app/solo /app/solo

ENTRYPOINT [ "/app/solo" ]