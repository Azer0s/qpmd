# Quacktor Port Mapper Daemon

## Prerequisites

* Go 1.15 or up

## Installation

qpmd can be install via `go get`/`go install`.

```
go get -u github.com/Azer0s/qpmd/...
go install github.com/Azer0s/qpmd/cmd/qpmd
```

or (as of Go 1.17):

```
go install github.com/Azer0s/qpmd/cmd/qpmd@latest
```

## Usage 

To get qmpd up and running you first have to install the daemon and then start it.
```
qpmd install
qpmd start
```

If you don't want qpmd as a daemon process and prefer to run it as a normal application instead, you can just run qpmd without any arguments.

```
qpmd
```

## Docker usage

The qpmd docker container uses [supervisord](http://supervisord.org/) under the hood. There already is a startup script for the quacktor app, so all you have to do is build your quacktor app and then place the binary into `/app/main`.

```Dockerfile
FROM golang:latest as build
WORKDIR /app
COPY . .
RUN go build github.com/my/app

FROM azer0s/qpmd:latest
COPY --from=build /app/main /app/main
```
