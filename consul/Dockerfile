FROM alpine:latest

# Install Consul
RUN apk add --no-cache curl unzip && \
    curl -o /tmp/consul.zip https://releases.hashicorp.com/consul/1.10.3/consul_1.10.3_linux_amd64.zip && \
    unzip /tmp/consul.zip -d /usr/local/bin/ && \
    rm /tmp/consul.zip

# Expose Consul ports
EXPOSE 8500 8600/udp

# Command to run Consul
CMD ["consul", "agent", "-dev", "-client=0.0.0.0"]