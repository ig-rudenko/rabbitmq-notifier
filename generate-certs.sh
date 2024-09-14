#!/bin/bash

# Проверяем количество параметров
if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <RabbitMQ domain> <Rabbit user>";
  exit 1;
fi

BASE_DIR="$(pwd)/rabbitmq-settings/";
HOST="$1";
RABBIT_USER="$2";

echo "$BASE_DIR";

docker run --rm \
  --mount type=bind,source="${BASE_DIR}",target=/rmq/ -w /rmq \
  -it ubuntu:24.04 /bin/bash -c "apt update && apt install openssl -y && ./create-certs.sh ${HOST} ${RABBIT_USER}";
