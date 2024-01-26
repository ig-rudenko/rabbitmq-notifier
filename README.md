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

### Файл конфигурации

Для работы приложения требуется настроить файл конфигурации `config.json`.

Через переменную окружения `CONFIG_FILE` можно указать где искать файл конфигурации,
по умолчанию - `./config.json`

## Запуск


### Отправка через producer

Для запуска **producer** необходимо обязательно указать в файле конфигурации все значения для блока `rabbitmq` и `exchange`.

```shell
notifier producer tg '{"chatId":123123, "message":"hello", "parseMode":"MarkdownV2", "token":"672364"}'
```

- параметр `producer` запускает приложение для отправки сообщения;
- `tg` это routingKey, который будет указан в сообщении;
- последний параметр это тело сообщения в формате JSON. 
Для обработки далее этого сообщения через `consumer` необходим JSON именно с такой структурой.


### Обработка через consumer

Для запуска **consumer** должны быть указаны все значения в файле конфигурации


```shell
notifier consumer telegram
```

- параметр `consumer` запускает приложения для приема сообщений;
- `telegram` это тип уведомителя, который должен обработать сообщение.
Доступны: `telegram`, `sms`, `email`.