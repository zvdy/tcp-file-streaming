# Use the official Nginx image as the base image
FROM nginx:latest

# Install Consul
RUN apt-get update && \
    apt-get install -y curl unzip && \
    curl -o /tmp/consul.zip https://releases.hashicorp.com/consul/1.10.3/consul_1.10.3_linux_amd64.zip && \
    unzip /tmp/consul.zip -d /usr/local/bin/ && \
    rm /tmp/consul.zip && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copy Nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Expose port 80
EXPOSE 80

# Command to run Consul and Nginx
CMD ["sh", "-c", "consul agent -dev -client=0.0.0.0 -join=consul & \
    sleep 10 && \
    curl --request PUT --data '{\"Name\": \"nginx\", \"Port\": 80, \"Tags\": [\"load-balancer\", \"nginx\", \"reverse-proxy\"]}' http://consul:8500/v1/agent/service/register && \
    nginx -g 'daemon off;'"]