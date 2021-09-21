FROM golang:latest as build
WORKDIR /qpmd
COPY . .
RUN go build cmd/qpmd/main.go

WORKDIR /pip
RUN apt-get update && apt-get install curl
RUN curl https://bootstrap.pypa.io/pip/2.7/get-pip.py --output get-pip.py

FROM ubuntu:latest
COPY --from=build /qpmd/main /app/qpmd
COPY --from=build /pip/get-pip.py get-pip.py

RUN apt-get update && apt-get install -y net-tools python2
RUN python2 get-pip.py
RUN mkdir -p /var/log/supervisor

RUN pip2 install supervisor
RUN pip2 install supervisor-stdout

COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY start_app.sh /app/start_app.sh
RUN chmod +x /app/start_app.sh
CMD ["/usr/local/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
