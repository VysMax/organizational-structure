# Организационная структура

Сервис для управления организационной структурой с поддержкой древовидных структур и рекурсивных зависимостей.

## Технологии

- Go 1.26.1
- PostgreSQL 17.4
- GORM
- Goose
- Docker / Docker Compose

### Запуск через Docker
```
docker-compose up -d --build
```

# Проверить, что все контейнеры запущены
```
docker-compose ps
```

# Остановка сервиса
```
docker-compose down
```

### API Эндпоинты
## 1. Создание департамента (CreateDepartment)
**POST /departments/**

Создаёт новый департамент.

# Ограничения:

- Название обязательно, длина от 1 до 200 символов

- Департамент не может быть родителем самого себя (проверяется на уровне БД)

# Запрос:

```json
{
    "name": "IT Department",
    "parent_id": 5
}

```

# Ответ:

201 Created

``` json
{
    "id": 1,
    "name": "IT Department",
    "parent_id": 5,
    "created_at": "2026-05-18T10:00:00Z"
}
```

## 2. Создание сотрудника (CreateEmployee)
**POST /departments/{id}/employees/**

Создаёт нового сотрудника в департаменте.

# Ограничения:

- Поле full_name обязательно (длина от 2 до 200 символов)

- Поле position обязательно

- Поле department_id обязательно и должно существовать

- Поле hired_at опционально, заполняется в формате DD-MM-YYYY

# Запрос:

``` json
{
    "full_name": "Mykola Honcharenko",
    "position": "Sales manager",
    "department_id": 1,
    "hired_at": "15-01-2026"
}
```

# Ответ (201 Created):

```json
{
    "id": 1,
    "department_id": 1,
    "full_name": "Mykola Honcharenko",
    "position": "Sales manager",
    "hired_at": "2026-01-15T00:00:00Z",
    "created_at": "2026-05-18T10:00:00Z"
}
```


## 3. Получение дерева департаментов сотрудниками (GetTree)
**GET /departments/{id}**

Возвращает департамент **id** со всеми дочерними подразделениями и, опционально, сотрудниками.

# Параметры:

depth — глубина рекурсии (1-5). Определяет, сколько уровней вложенности включать в ответ, **не считая целевой департамент**.

include_employees — включать ли сотрудников в ответ

# Ограничения:

- Сотрудники сортируются по full_name по алфавиту

- Максимальная depth — 5 уровней

# Пример запроса:

```json
{
    "depth":5,
    "include_employees":true
}

```

# Ответ

```json
{
    "id": 4,
    "name": "International Cooperation",
    "parent_id": 2,
    "created_at": "2026-05-18T15:23:14.124663Z",
    "employees": [
        {
            "id": 4,
            "departament_id": 4,
            "full_name": "Alesia",
            "position": "Boss",
            "created_at": "2026-05-18T15:52:50.552807Z"
        },
        {
            "id": 2,
            "departament_id": 4,
            "full_name": "Razvan",
            "position": "Assistant",
            "hired_at": "2020-02-01T00:00:00Z",
            "created_at": "2026-05-18T15:26:18.774927Z"
        }
    ],
    "children": [
        {
            "id": 11,
            "name": "Europe",
            "parent_id": 4,
            "created_at": "2026-05-18T16:44:36.462172Z"
        },
        {
            "id": 12,
            "name": "Asia",
            "parent_id": 4,
            "created_at": "2026-05-18T18:08:55.900066Z",
            "employees": [
                {
                    "id": 9,
                    "departament_id": 12,
                    "full_name": "Rahim Mammadov",
                    "position": "Head",
                    "created_at": "2026-05-18T21:34:06.531201Z"
                },
                {
                    "id": 12,
                    "departament_id": 12,
                    "full_name": "Ihar Piatrenka",
                    "position": "Assistant",
                    "created_at": "2026-05-18T21:47:14.682397Z"
                }
            ]
        },
        {
            "id": 5,
            "name": "Africa",
            "parent_id": 4,
            "created_at": "2026-05-18T15:23:24.319041Z",
            "children": [
                {
                    "id": 19,
                    "name": "Subsaharian Africa",
                    "parent_id": 5,
                    "created_at": "2026-05-18T21:33:23.763457Z"
                }
            ]
        }
    ]
}
```

## 4. Смена родителя департамента (UpdateParent)
**PATCH /departments/{id}**

Переподчиняет департамент **id** другому департаменту и, опционально, переименовывает.

# Ограничения:

- Запрещено переподчинение департаменту подразделению, находящемуся внутри его поддерева (предотвращение создания циклов в дереве), в том числе самому себе

- Новый родитель должен существовать

- parent_id = null делает департамент корневым

# Запрос:

```json
{
    "name": "Innovation department",
    "parent_id": 19
}
```
# Ответ 
```json
{
    "id": 23,
    "name": "Innovation department",
    "parent_id": 19,
    "created_at": "2026-05-18T22:24:28.579837Z"
}
```

## 5. Удаление департамента (DeleteDepartment)
**DELETE /departments/{id}**

Удаляет департамент **id**. Предусмотрено два режима обработки зависимостей.

# Параметры:

- mode - режим удаления:

"cascade" - каскадное удаление (удаляются все дочерние департаменты и сотрудники). Реализовано на уровне базы данных

"reassign" - переподчинение (дочерние департаменты переподчиняются родителю удаляемого, сотрудники переводятся в обязательно отдельно указываемый департамент)

- reassign_to_department_id — ID департамента, куда перевести сотрудников (обязателен для mode=reassign)

# Ограничения:

- При mode="reassign":

reassign_to_department_id должен существовать

# Пример запроса:


- Каскадное удаление
```json
{
    "mode":"reassign",
    "reassign_to_department_id": 12
}
```

- Удаление с переподчинением
```json
{
    "mode":"reassign",
    "reassign_to_department_id": 12
}
```

# Ответ: 204 No Content

### Тестирование

# Юнит-тесты
```
go test ./usecase/...
```