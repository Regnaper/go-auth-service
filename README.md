
# Golang cервис авторизации 

> Приложение выдает клиенту JWT access и refresh токены.

> При наличии у клиента cookie с refresh токеном приложение выполняет для него refresh операцию.

> Токены хранятся в виде JSON в MongoDB.

**Маршруты**

- secret/{guid} (method: GET)
-- для пользователя guid создает в базе связанные access и refresh токены и устанавливает их в cookies клиента;
- secret/{guid} (method: POST)
-- получает из cookies клиента refresh токен и для пользователя guid обновляет (при наличии) этот refresh токен 
и связанный с ним access токен, после чего передает их в cookies клиента;
- secret/{guid} (method: DELETE)
-- получает из cookies клиента refresh токен и для пользователя guid удаляет (при наличии) этот refresh токен 
   и связанный с ним access токен;
- secrets/{guid} (method: DELETE)
-- получает из cookies клиента refresh токен и удаляет для пользователя guid все refresh токены и связанные с ними access токены.

Результат: https://rocky-escarpment-09841.herokuapp.com/