#!/bin/bash

# SSL Certificate Generation Script for Domain MAX
# This script generates self-signed certificates for development/testing
# For production, replace with certificates from a trusted CA

set -e

# Configuration
SSL_DIR="./ssl"
CERT_FILE="$SSL_DIR/nginx-selfsigned.crt"
KEY_FILE="$SSL_DIR/nginx-selfsigned.key"
COUNTRY="US"
STATE="State"
LOCALITY="City"
ORGANIZATION="Domain MAX"
ORGANIZATIONAL_UNIT="IT Department"
COMMON_NAME="localhost"
EMAIL="admin@domain-max.com"

echo "🔐 Generating SSL certificates for Domain MAX..."

# Create SSL directory
mkdir -p "$SSL_DIR"

# Generate private key
echo "📝 Generating private key..."
openssl genrsa -out "$KEY_FILE" 4096

# Generate certificate signing request
echo "📄 Generating certificate signing request..."
openssl req -new -key "$KEY_FILE" \
    -out "$SSL_DIR/nginx-selfsigned.csr" \
    -subj "/C=$COUNTRY/ST=$STATE/L=$LOCALITY/O=$ORGANIZATION/OU=$ORGANIZATIONAL_UNIT/CN=$COMMON_NAME/emailAddress=$EMAIL"

# Generate self-signed certificate
echo "🏆 Generating self-signed certificate..."
openssl x509 -req -days 365 \
    -in "$SSL_DIR/nginx-selfsigned.csr" \
    -signkey "$KEY_FILE" \
    -out "$CERT_FILE" \
    -extensions v3_req \
    -extfile <(cat <<EOF
[v3_req]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = domain-max.local
DNS.3 = *.domain-max.local
IP.1 = 127.0.0.1
IP.2 = ::1
EOF
)

# Set proper permissions
chmod 600 "$KEY_FILE"
chmod 644 "$CERT_FILE"

# Clean up CSR file
rm -f "$SSL_DIR/nginx-selfsigned.csr"

echo "✅ SSL certificates generated successfully!"
echo "📁 Certificate: $CERT_FILE"
echo "🔑 Private key: $KEY_FILE"
echo ""
echo "⚠️  Note: These are self-signed certificates for development/testing only."
echo "📚 For production, obtain certificates from a trusted Certificate Authority."
echo ""
echo "🚀 You can now start the services with: docker-compose up -d"