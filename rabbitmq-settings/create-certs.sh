#!/bin/bash

HOST="$1";
RABBIT_USER="$2";

echo "RabbitMQ domain = $HOST";
echo "RabbitMQ user = $RABBIT_USER";

# Создаем конфигурационный файл для OpenSSL с SAN
cat > openssl.cnf <<EOL
[ req ]
default_bits = 2048
prompt = no
default_md = sha256
distinguished_name = dn
req_extensions = req_ext

[ dn ]
C = RU
ST = SEV
L = SEV
O = SC
OU = NOC
emailAddress = example@mail.com
CN = $HOST

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = $HOST

[ v3_ca ]
subjectAltName = @alt_names
EOL

#### Create a Certificate Signing Request (CSR) with SAN
openssl req -days 3650 -new -text -nodes -keyout server.key -out server.csr -config openssl.cnf;

#### Generate self-signed certificate with SAN
openssl x509 -days 3650 -req -in server.csr -signkey server.key -out server.crt -extensions req_ext -extfile openssl.cnf;

#### Also make the server certificate the root-CA certificate
cp server.crt root.crt;

#### Remove the now-redundant CSR
rm server.csr;

### Generate client certificates to be used by clients/connections

#### Create a Certificate Signing Request (CSR) for the client
openssl req -days 3650 -new -nodes -subj "/C=RU/ST=SEV/L=SEV/O=SC/OU=NOC/emailAddress=example@mail.com/CN=$RABBIT_USER" -keyout client.key -out client.csr;

#### Create a signed certificate for the client using our root certificate
openssl x509 -days 3650 -req -CAcreateserial -in client.csr -CA root.crt -CAkey server.key -out client.crt;

#### Remove the now-redundant CSR
rm client.csr;

# Удаляем временный файл конфигурации
rm openssl.cnf;
