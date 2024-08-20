FROM golang:latest
MAINTAINER Tonywww

RUN mkdir -p /work

WORKDIR /work

COPY . /work

RUN go build /monibuca/main.go

EXPOSE 27092

