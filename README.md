# URL Shortener API

An educational backend service for creating short URLs, built with Go.

## Tech Stack

- Go, Chi
- PostgreSQL, pgxpool
- Docker Compose
- SQL migrations

## Getting Started

Create an environment file:

```bash
cp .env.example .env
```

Set the PostgreSQL user, password, and database name in `.env`, then start PostgreSQL:

```bash
docker compose up -d db
```

Apply the migrations:

```bash
migrate \
  -path ./migrations \
  -database "postgres://postgres:postgres@localhost:5432/shortener?sslmode=disable" \
  up
```

The connection values in this command must match your `.env` file.

Build and start the application:

```bash
docker compose up -d --build app
```

Check that the service is running:

```bash
curl http://localhost:8080/health
```

## API

- `POST /api/v1/urls` — create a short URL
- `GET /api/v1/urls/{shortCode}` — get URL information
- `GET /{shortCode}` — redirect to the original URL
- `DELETE /api/v1/urls/{shortCode}` — delete a short URL
- `GET /health` — check service health

Create a short URL:

```bash
curl -X POST http://localhost:8080/api/v1/urls \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}'
```

## Tests

With the test database running:

```bash
DATABASE_URL="postgres://postgres:postgres@localhost:5432/shortener?sslmode=disable" go test ./...
```

## Stopping the Application

```bash
docker compose down
```

PostgreSQL data is persisted in the `db-data` Docker volume.
