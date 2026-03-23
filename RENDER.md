# Render Setup

This service is intended to run as a private Render service behind the public gateway.

## What changed for Render

- The service now supports Render's `PORT` environment variable automatically.
- The database config now supports `DB_URL` in addition to `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, and `DB_NAME`.

## Recommended deployment

1. Create a Render Postgres database.
2. Deploy `curriculum-service` as a private service.
3. Set `DB_URL` from the database `connectionString`.
4. Apply SQL migrations from the `migrations/` directory.
5. After the service is healthy, connect the public gateway to this private service by its Render internal `hostport`.

## Health check

- `/health`
