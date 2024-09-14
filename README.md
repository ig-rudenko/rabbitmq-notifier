# RMQ Notifier

Утилита для управления очередями сообщений с использованием RabbitMQ

Данный инструмент предназначен для запуска Producer'ов и Consumer'ов в системе RabbitMQ. Он эффективно обрабатывает сообщения, обеспечивая их постепенную передачу для избежания перегрузки и зависания приложений.

Основные возможности:

* Consumer поддерживает отправку уведомлений в Telegram и на Email.
* Обеспечивает плавную передачу уведомлений без резких всплесков нагрузки.
* Идеально подходит для асинхронной отправки уведомлений и писем, что помогает улучшить производительность и стабильность ваших приложений.
* Producer можно использовать для передачи любых сообщений для последующей обработки другими Consumers.

![Basic Usage.png](docs/img/BasicUsage.png)

## Схема работы producer-consumer

![schema.png](docs/img/schema.png)

## Схема работы only-producer

Producer не проверяет тело сообщения, так что можно передавать любые данные

![schema.png](docs/img/only-producer.png)

## Быстрый запуск через docker compose

Будут подняты следующие контейнеры:

- RabbitMQ
- Producer
- Telegram Consumer
- Email Consumer

1. Создать сертификаты

    ```shell
    sudo bash generate-certs.sh 'rabbitmq' 'rmuser';
    ```

2. Запустить контейнеры

    ```shell
    sudo docker compose up -d;
    ```

3. Передать POST запрос на URL `http://localhost:9090/<routingKey>`. Пример для телеграма:

    ```bash
    curl -X POST "http://localhost:9090/telegram" \
    -H "Authorization: Token insecureToken" \
    -H "Content-Type: application/json" \
    -d '{
      "chatId": 123123123,
      "message": "hello",
      "parseMode": "MarkdownV2",
      "token": "YOUR_TELEGRAM_TOKEN"
    }'
    ```


## 1. Настройка

### 1.1. Создание сертификатов

Создать файлы сертификатов RabbitMQ для сервера и клиента можно через скрипт `generate-certs.sh`

```shell
sudo bash generate-certs.sh 'rabbitHostDomain' 'rabbitUser';
```

Чтобы сгенерировать сертификаты для работы `docker compose`, нужно указать:

```shell
sudo bash generate-certs.sh 'rabbitmq' 'rmuser';
```

Скрипт создаст 5 файлов, которые будут в папке `rabbitmq-settings`:

<details>

- `server.key`: Этот файл содержит закрытый ключ (private key) сервера.
  Этот ключ используется для подписания запроса на сертификат и для создания подписанного
  самоподписанного серверного сертификата.
- `server.crt`: Самоподписанный серверный сертификат. Этот файл содержит открытый ключ сервера,
  данные о сервере и подпись, сгенерированную закрытым ключом сервера. Он используется
  сервером для установки защищенного соединения.
- `root.crt`: Корневой сертификат (cacert). В данном случае, этот файл представляет
  собой самоподписанный корневой сертификат, который используется для подписи клиентских сертификатов.
  Этот файл может быть распространен среди клиентов для проверки подлинности сервера.
- `client.key`: Закрытый ключ (private key) клиента. Этот ключ используется для создания
  запроса на подпись сертификата клиента и для создания подписанного клиентского сертификата.
- `client.crt`: Подписанный клиентский сертификат. Этот файл содержит открытый ключ клиента,
  данные о клиенте и подпись, сгенерированную закрытым ключом сервера (root key).

</details>


Для запуска rabbitmq будут необходимы `server.key`, `server.crt`, `root.crt`.

Для подключения клиента - `root.crt`, `client.key`, `client.crt`.

### 1.2. Файл конфигурации

Для работы приложения требуется настроить файл конфигурации `config.json`.

Через переменную окружения `CONFIG_FILE` можно указать, где искать файл конфигурации,
по умолчанию - `/etc/rmq-notifier/config.json`

Каждое значение файла конфигурации можно переопределить через переменную окружения. 

Переменные окружения:

```shell
RABBITMQ_USER: # Имя пользователя для подключения к RabbitMQ.
RABBITMQ_PASSWORD: # Пароль для подключения к RabbitMQ.
RABBITMQ_HOST: # Хост RabbitMQ сервера.
RABBITMQ_PORT: # Порт RabbitMQ сервера.
RABBITMQ_VHOST: # Виртуальный хост RabbitMQ, который будет использоваться.
RABBITMQ_CACERT: # Путь к файлу корневого сертификата для проверки TLS соединений RabbitMQ.
RABBITMQ_CERTFILE: # Путь к файлу сертификата клиента для TLS соединений RabbitMQ.
RABBITMQ_KEYFILE: # Путь к файлу ключа клиента для TLS соединений RabbitMQ.

EXCHANGE_NAME: # Имя обмена RabbitMQ, которое будет использоваться для публикации/подписки на сообщения.
EXCHANGE_TYPE: # Тип обмена RabbitMQ (например, direct, topic, fanout).

CONSUMER_CONNECTION_NAME: # Имя подключения для consumer.
CONSUMER_ROUTING_KEY: # Ключ маршрутизации для consumer.
CONSUMER_QUEUE: # Очередь RabbitMQ, которую будет использовать consumer.
CONSUMER_COUNT: # Количество экземпляров consumer, которые будут запущены.

# Количество сообщений, которые consumer может забрать из очереди перед подтверждением.
CONSUMER_PREFETCH_COUNT:

# Время в секундах, через которое сообщения будут помечены как Acknowledge и пропущены, если они не были обработаны.
CONSUMER_EXPIRE_AFTER_SECONDS:

PRODUCER_AUTH_TOKEN: # Токен авторизации для producer.

# Для email consumer (подключение к почтовому серверу)
EMAIL_NOTIFIER_HOST: Хост почтового сервера для отправки уведомлений по email.
EMAIL_NOTIFIER_PORT: Порт почтового сервера.
EMAIL_NOTIFIER_LOGIN: Логин для подключения к почтовому серверу.
EMAIL_NOTIFIER_PASSWORD: Пароль для подключения к почтовому серверу.
```

## 2. Запуск

### 2.1. Отправка через producer

Запускаем приложение на конкретном порту

```shell
notifier producer 0.0.0.0:9090
```

Далее отправляем POST запрос на URL `/<routingKey>`.

```bash
curl -X POST "http://localhost:9090/telegram" \
-H "Authorization: Token ********" \
-H "Content-Type: application/json" \
-d '{
  "chatId": 123123123,
  "message": "hello",
  "parseMode": "MarkdownV2",
  "token": "****"
}'
```

Каждый запрос должен содержать заголовок с токеном, который указан в файле конфигурации,
либо через переменную окружения `PRODUCER_AUTH_TOKEN`
```json
"producer": {
  "authToken": "834932789472389478923"
}
```

### 2.2. Обработка через consumer


```shell
notifier consumer telegram
```

- параметр `consumer` запускает приложения для приема сообщений;
- `telegram` это тип уведомителя, который должен обработать сообщение.
Доступны: `telegram`, `email`.

### 2.2.1. Telegram

Для уведомителя `telegram` тело сообщения должно быть в следующем JSON в формате:

```json
{ 
  "chatId": 123123123,
  "message": "hello", 
  "parseMode": "MarkdownV2",
  "token": "****"
}
```

### 2.2.2. Email

Для уведомителя `email` тело сообщения должно быть в следующем JSON в формате:

```json
{
  "sender": "user@mail.com",
  "to": [
    "to-user1@mail.com",
    "to-user2@mail.com"
  ],
  "subject": "test",
  "body": "<h1>test</h1>"
}
```

> [!IMPORTANT]
> Для работы `email` уведомителя необходимо в файле конфигураций указать настройки для подключения,
> либо через переменные окружения:

```json
{
  "emailNotifier": {
    "host": "mail.domain",
    "port": 587,
    "login": "user",
    "password": "password"
  }
}
```
