# Восемь тестовых запросов

PowerShell, сервер запущен на порту 8080. POST создаёт первую задачу в новом процессе.

1. Состояние сервера → 200, {"status":"ok"}.
```powershell
curl.exe -i http://localhost:8080/health
```
2. Создание задачи → 201, id=1, done=false.
```powershell
Invoke-WebRequest -Method Post -Uri http://localhost:8080/tasks -ContentType 'application/json' -Body '{"title":"Купить молоко"}' | Select-Object StatusCode,Headers,Content
```
3. Список задач → 200, массив с задачей.
```powershell
curl.exe -i http://localhost:8080/tasks
```
4. Поиск без учёта регистра → 200, массив с задачей.
```powershell
curl.exe -i 'http://localhost:8080/tasks?q=%D0%9C%D0%9E%D0%9B%D0%9E%D0%9A%D0%9E'
```
5. Получение задачи → 200, объект с id=1.
```powershell
curl.exe -i http://localhost:8080/tasks/1
```
6. Изменение done → 200, done=true.
```powershell
Invoke-WebRequest -Method Patch -Uri http://localhost:8080/tasks/1 -ContentType 'application/json' -Body '{"done":true}' | Select-Object StatusCode,Headers,Content
```
7. Удаление → 204, пустое тело.
```powershell
curl.exe -i -X DELETE http://localhost:8080/tasks/1
```
8. Ошибка ID → 400, JSON-ошибка invalid id.
```powershell
curl.exe -i http://localhost:8080/tasks/abc
```

Проверки 404, некорректного JSON, длины title, CORS и конфигурации порта приведены в testing.md.
