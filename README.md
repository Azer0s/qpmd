# Quacktor Port Mapper Daemon

## Prerequisites

* Go 1.15 or up

## Installation

qpmd can be install via `go get`/`go install`.

```
go get -u  github.com/Azer0s/qpmd/...
go install github.com/Azer0s/qpmd/cmd/qpmd
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
