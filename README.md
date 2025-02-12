#                               **gw REST Api**

### Для запуска приложения:
Сначала скачайте проект к себе в окружение:
```
git clone https://github.com/Njrctr/gw-currency-wallet && cd gw-currency-wallet
```
Запуск:
```
make run
```

Если приложение запускается впервые, необходимо применить миграции к базе данных:

```
make migrate.up
```

### В данном проекте был реализован REST Api сервис для работы с кошельками:
* Весь функционал можно протестировать в SWAGGER документации: http://localhost:8080/swagger/
* Реализован функционал Аутентификации на основе JWT токена
* Реализован ИнМемори кэш
* REST сервис общается по gRPC с сервисом gw-exchanger https://github.com/Njrctr/gw-exchanger
* Контракт proto https://github.com/Njrctr/gw-proto-exchange