#!/bin/bash

# Wait until qpmd started and listens on port 7161
# shellcheck disable=SC2006
while [ -z "`netstat -tln | grep 7161`" ]; do
  echo 'Waiting for qpmd to start'
done
echo 'qpmd started'

# Start quacktor app
./app/main
