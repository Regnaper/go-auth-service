
# Golang cервис авторизации 

> Приложение выдает клиенту JWT access токен и refresh токен.

> При наличии у клиента cookie с refresh токеном приложение выполняет для него refresh операцию.

> Refresh токены хранятся в виде JSON в MongoDB.

**Маршруты**

- secret/{guid} (method: GET)
-- для пользователя guid создает access и refresh токены и устанавливает их в cookies клиента;
- secret/{guid} (method: POST)
-- получает из cookies клиента refresh токен и для пользователя guid обновляет (при наличии) этот refresh токен 
и связанный с ним access токен, после чего передает их в cookies клиента;
- secret/{guid} (method: DELETE)
-- получает из cookies клиента refresh токен и для пользователя guid удаляет (при наличии) этот refresh токен;
- secrets/{guid} (method: DELETE)
-- получает из cookies клиента refresh токен и (при его наличии в базе) удаляет для пользователя guid все refresh токены.

Результат: https://rocky-escarpment-09841.herokuapp.com/