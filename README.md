# Практическая работа 3 — HTTP-сервер на net/http

Сортов Семён Александрович, ЭФМО-01-26. Дата: 09.10.2026.

HTTP-сервис задач на стандартной библиотеке net/http. Поддерживает список и поиск по q, создание, получение, изменение done и удаление. Дополнения: валидация title 3–140 символов, CORS, логирование статуса и времени, graceful shutdown; домашнее задание: PORT, tasks.ps1 и два юнит-теста обработчиков. Требования: Go 1.22+ и Git.

## Запуск и сборка

Команды выполняются из этой папки. Каждая работа запускается отдельно: по умолчанию все используют порт 8080.

```powershell
go mod download
go run ./cmd/server
# Сборка и проверки
go build -o ./bin/server.exe ./cmd/server
go fmt ./...
go vet ./...
go test ./...
# После остановки go run
.\bin\server.exe
```

## Примеры запросов

```powershell
curl.exe -i http://localhost:8080/health
curl.exe -i http://localhost:8080/tasks
# Восемь примеров запросов находятся в requests.md
```

## Конфигурация

```powershell
$env:PORT="8082"
.\tasks.ps1 run
.\tasks.ps1 build
.\tasks.ps1 test
```
Без PORT используется 8080. Для завершения сервера нажмите Ctrl+C. Данные хранятся в памяти и очищаются при перезапуске.

## Структура

```text
.gitignore
cmd/server/main.go
go.mod
internal/api/handlers.go
internal/api/handlers_test.go
internal/api/middleware.go
internal/api/responses.go
internal/storage/memory.go
tasks.ps1
testing.md
bin/server.exe (собранная программа, исключена из Git)
```
