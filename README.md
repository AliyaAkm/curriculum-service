# curriculum-service

Go-сервис для учебного каталога с поиском и фильтрацией курсов под фронтенд образовательного портала.

## Что уже реализовано

- поиск по `q` по локализованным названиям, описаниям, темам и тегам;
- фильтрация по темам, уровню, минимальному рейтингу, длительности и наличию сертификата;
- мультиязычная выдача `ru/en/kz`;
- API для получения доступных фильтров;
- Postgres-схема и демо-данные.

## Структура API

### `GET /health`

Проверка состояния сервиса.

### `GET /courses`

Возвращает общий список всех опубликованных курсов без текстового поиска и без структурной фильтрации.

Параметры:

- `locale` - `ru`, `en`, `kz`;
- `page` - номер страницы;
- `page_size` - размер страницы.

Пример:

```bash
curl "http://localhost:8083/courses?locale=en&page=1&page_size=12"
```

### `GET /courses/search`

Параметры:

- `q` - поисковая строка;
- `locale` - `ru`, `en`, `kz`;
- `page` - номер страницы;
- `page_size` - размер страницы.

Пример:

```bash
curl "http://localhost:8083/courses/search?q=qa&locale=en&page=1&page_size=12"
```

### `GET /courses/filter`

Параметры:

- `locale` - `ru`, `en`, `kz`;
- `topic` или `topics` - одно или несколько значений;
- `level` или `levels` - `beginner`, `intermediate`, `advanced`;
- `min_rating` - число от `0` до `5`;
- `duration` или `durations` - `quick`, `focused`, `deep`;
- `with_certificate` - `true`;
- `page` - номер страницы;
- `page_size` - размер страницы.

Пример:

```bash
curl "http://localhost:8083/courses/filter?locale=en&topic=programming-languages&level=intermediate&min_rating=4&duration=quick&with_certificate=true&page=1&page_size=12"
```

### `GET /courses/filters`

Возвращает темы, уровни, длительности, рейтинги и наличие сертификатов для построения фильтров на фронтенде.

Пример:

```bash
curl "http://localhost:8083/courses/filters?locale=ru"
```

## Запуск

1. Скопируйте `.env.example` в `.env`.
2. Примените миграции:

```bash
psql -d curriculum -f migrations/001_init.sql
psql -d curriculum -f migrations/002_seed_demo.sql
```

3. Запустите сервис:

```bash
go run ./cmd/curriculum
```

## Архитектура

- `internal/domain` - доменные типы и валидация;
- `internal/usecase` - бизнес-логика каталога;
- `internal/repo/postgres` - SQL-репозиторий;
- `internal/http` - DTO, handlers, router;
- `cmd/curriculum` - точка входа.
