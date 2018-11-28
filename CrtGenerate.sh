#!/bin/bash

# create self-signed server certificate:

read -p "Enter your domain: " DOMAIN

echo "Create server key..."

openssl genrsa -des3 -out $DOMAIN.key 2048

echo "Create server certificate signing request..."

#SUBJECT="/C=US/ST=$DOMAIN/L=$DOMAIN/O=$DOMAIN/OU=$DOMAIN/CN=$DOMAIN"
SUBJECT="/CN=$DOMAIN"

openssl req -new -subj $SUBJECT -key $DOMAIN.key -out $DOMAIN.csr

echo "Remove password..."

mv $DOMAIN.key $DOMAIN.origin.key
openssl rsa -in $DOMAIN.origin.key -out $DOMAIN.key

echo "Sign SSL certificate..."

openssl x509 -req -days 3650 -in $DOMAIN.csr -signkey $DOMAIN.key -out $DOMAIN.crt

echo "Generated:"
echo "    $DOMAIN.crt"
echo "    $DOMAIN.key"
