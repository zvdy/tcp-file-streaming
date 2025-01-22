#!/bin/bash

# Define the directory to store the certificate and key
CERT_DIR="./nginx/tls"
mkdir -p $CERT_DIR

# Generate a self-signed certificate and key
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout $CERT_DIR/nginx.key -out $CERT_DIR/nginx.crt -subj "/C=US/ST=State/L=City/O=Organization/OU=Department/CN=localhost"

echo "Self-signed certificate and key generated at $CERT_DIR"