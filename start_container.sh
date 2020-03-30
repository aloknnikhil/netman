#!/bin/bash

# Stop existing containers
docker kill agent1 || true
docker kill etcd || true

# Start ETCD
docker run -dit --rm -p 2379:2379 --name etcd --rm quay.io/coreos/etcd:v3.1.0 /usr/local/bin/etcd -advertise-client-urls http://0.0.0.0:2379 -listen-client-urls http://0.0.0.0:2379

# Start VPP
docker run -dit --rm --name agent1 -p 9111:9111 -e MICROSERVICE_LABEL="meter_vpp" --privileged ligato/vpp-agent
