# gophermart

Накопительная система лояльности «Гофермарт». Индивидуальный дипломный проект курса «Go-разработчик»

В качестве фреймворка для обработки запросов выбран [echo](https://github.com/labstack/echo).

## План реализации

### 1. Структура проекта

- [x] настройка автотестов
- [x] docker compose для локальной разработки
- [x] выбор фреймворков и библиотек (echo, sqlx)
  - [x] выбор логгера (slog)
- [ ] линтеры и форматтеры: `golangci-lint`, `goimports`, `gofmt`

### 2. Основные модели данных и миграции

- [ ] `users` – пользователи
- [ ] `orders` – заказы (номера)
- [ ] `transactions` – транзакции (пополнения и списания)

### 3. Регистрация, аутентификация и авторизация пользователей

- [ ] `POST /api/user/register` — регистрация пользователя
- [ ] `POST /api/user/login` — аутентификация пользователя
- Настройка приватного ключа
- Middleware для авторизации запросов

### 4. Работа с заказами

- [ ] `POST /api/user/orders` — загрузка пользователем номера заказа для расчёта 
  - [ ] регистрация заказа и привязка к пользователю
- [ ] `GET /api/user/orders` — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях

### 5. Взаимодействие с системой расчета баллов лояльности

- [ ] Проверка заказа в системе accrual и начисление баллов (поллинг, воркер пул)

### 6. Баланс

- [ ] `GET /api/user/balance` — получение текущего баланса счёта баллов лояльности пользователя

### 7. Начисление и списание баллов, получение истории списаний

- [ ] `POST /api/user/balance/withdraw` — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
- [ ] `GET /api/user/withdrawals` — получение информации о выводе средств с накопительного счёта пользователем

### 8. Тестирование

- [ ] юнит-тесты
- [ ] тесты API
- [ ] интеграционные тесты

### 9. Документация

- [x] `README.md` с описанием проекта и планом реализации
- [ ] API документация (swagger)

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m master template https://github.com/yandex-praktikum/go-musthave-diploma-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/master .github
```

Затем добавьте полученные изменения в свой репозиторий.
