FROM golang:tip-alpine3.23

LABEL maintainer="maksimacx50@gmail.com"

COPY . /one-team/app/source

WORKDIR /one-team/app/source

RUN go mod tidy && go build ./cmd/app

WORKDIR /one-team/app/run

RUN mv /one-team/app/source/app . && \ 
    mv /one-team/app/source/config.yaml . && \ 
    mkdir /one-team/app/run/logs

VOLUME /one-team/app/run/logs

EXPOSE 8080

ENTRYPOINT ["./app"]
