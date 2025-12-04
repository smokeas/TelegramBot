Типы запросов
Get - получить данные
Post - отправить данные 
Patch / put - обновить данные 
DElete - удалить данные 

# Получить значение
curl http://localhost:8080/get

# Увеличить счётчик
curl -X POST http://localhost:8080/post

 # Проверка 
 
 # Сначала GET
curl.exe http://localhost:8080/get

# Потом POST
curl.exe -X POST http://localhost:8080/post

# И снова GET — должен быть counter = 1
curl.exe http://localhost:8080/get

для работы с json
 go get github.com/labstack/echo/v4  