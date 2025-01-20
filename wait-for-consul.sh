#!/bin/sh

# Wait for Consul to be ready
while ! curl -s http://localhost:8500/v1/status/leader > /dev/null; do
  echo "Waiting for Consul to be ready..."
  sleep 1
done

echo "Consul is ready"