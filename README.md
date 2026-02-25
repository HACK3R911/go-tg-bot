# Telegram bot

Это самый простой телеграм бот который по команде /snake отправляет ссылку на последнее видео из ютуба

## (TODO)
1. создать новую ветку из dev 
2. настроить cicd пайплайн с линтом/тестами/пушем в докер регистр
3. сделать новую ветку из dev для создания rate limit сообщений
4. реализовать rate limit
5. реализация кеша
6. di контейнер

## 🔴 Критические проблемы

[//]: # (1. bot.go:69 - Сообщение создаётся, но не отправляется:)
[//]: # (   default:)
[//]: # (   tgbotapi.NewMessage... // не отправляется!)
[//]: # ()
[//]: # (2. cmd.go:46 - Аналогичная проблема: bot.Send без проверки ошибки)
3. Нет rate limiting - Каждый message запускает горутину без ограничений (OOM/бан от Telegram)
4. cmd.go:40 - context.Background() вместо переданного контекста

## ⚠️ Высокий приоритет
1. cache.go пустой - YouTube API имеет квоты, нужен кэш (Redis)
2. Нет структурированного логирования - используется log вместо zerolog/zap
3. main.go:55 - sql.Open не проверяет соединение перед миграциями
4. Graceful shutdown - горутины не отслеживаются, возможны утечки
5. Нет метрик - нет Prometheus/OpenTelemetry
6. Нет retry логики - для внешних API (Telegram, YouTube)

## 📌 Средний приоритет
1. Тесты - только auth repository покрыт, нет handler/service тестов
2. Health check endpoint - нужен для production
3. Нет Docker - отсутствует Dockerfile
4. Circuit breaker - для внешних API
5. Config validation - порты, URL и т.д. не валидируются
6. Rate limiting per user - спам от одного пользователя

##  . Низкий приоритет
1. Hardcoded strings (команды, сообщения)
2. Нет docker-compose для dev
3. Нет graceful degradation при ошибках API
4. DB connection pool не настроен (pgxpool конфиг)