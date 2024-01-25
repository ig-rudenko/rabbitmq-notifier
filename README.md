## Настройка

### Создание сертификатов


Для работы библиотеки x509 необходимо настроить `subjectAltName`. Для этого добавьте в файл `/etc/ssl/openssl.conf`
блок `alt_names` где будут альтернативные домены. Обязательно укажите основной домен, так как subjectAltName
переопределяет полностью CN.

```ini
[ alt_names ]
DNS.1 = server1.example.com
DNS.2 = mail.example.com
DNS.3 = www.example.com
DNS.4 = www.sub.example.com
DNS.5 = mx.example.com
DNS.6 = support.example.com

[ v3-ca ]
# Подключаем
subjectAltName = @alt_names
...
```

Создать файлы сертификатов для сервера и клиента можно через скрипт `rabbit-settings/create-certs.sh`

```shell
bash rabbit-settings/create-certs.sh 'rabbitHost' 'rabbitUser';
```

### Переменные окружения

Для работы приложения требуются перменные окружения.
Необязательные переменные уже имеют значения.


    CACERT         # root cert
    CERTFILE       # client cert
    KEYFILE        # client key

    RABBITMQ_USER
    RABBITMQ_PASS
    RABBITMQ_HOST
    RABBITMQ_PORT="5671"
    RABBITMQ_VHOST=""

    RABBITMQ_CONNECTION_NAME            # Название подписчика
    RABBITMQ_EXCHANGE_NAME              # Название точки обмена
    RABBITMQ_EXCHANGE_TYPE="direct"     # Тип обменника
    RABBITMQ_ROUTING_KEY                # Ключ маршрутизации
    RABBITMQ_QUEUE                      # Очередь для подключения
    CONSUMER_COUNT="3"
    PREFETCH_COUNT="1"
