#!/bin/bash

HOST="$1";
RABBIT_USER="$2";

echo "RabbitMQ domain = $HOST";
echo "RabbitMQ user = $RABBIT_USER";

#### Create a Certificate Signing Request (CN=Server domain)
openssl req -days 3650 -new -text -nodes -subj "/C=RU/ST=SEV/L=SEV/O=SC/OU=NOC/emailAddress=example@mail.com/CN=$HOST" -keyout server.key -out server.csr;

#### Generate self-signed certificate
openssl req -days 3650 -x509 -text -in server.csr -key server.key -out server.crt;

#### Also make the server certificate to be the root-CA certificate
cp server.crt root.crt;

#### Remove the now-redundant CSR
rm server.csr;


### Generate client certificates to be used by clients/connections

#### Create a Certificate Signing Request (CN=RabbitMQ user)
openssl req -days 3650 -new -nodes -subj "/C=RU/ST=SEV/L=SEV/O=SC/OU=NOC/emailAddress=example@mail.com/CN=$RABBIT_USER" -keyout client.key -out client.csr;

#### Create a signed certificate for the client using our root certificate
openssl x509 -days 3650 -req -CAcreateserial -in client.csr -CA root.crt -CAkey server.key -out client.crt;

#### Remove the now-redundant CSR
rm client.csr;

