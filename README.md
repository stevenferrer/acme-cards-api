# ACME Cards API

An example of card-as-service integration with [Reap API](https://reap.readme.io/reference/test-environment) sandbox environment.

See also [ACME Cards Dashboard](https://github.com/stevenferrer/acme-cards-dashboard).

## Requirements

- [Go](https://go.dev/dl) 1.25 or newer
- [PostgreSQL](https://www.postgresql.org) 18 or newer

## Guides

### Export .env

```sh
export $(grep -v '^#' .env | xargs)
```

### Database migration

Export env variable `POSTGRES_DSN`, refer to [.env.example](.env.example) file.

Build and run the `migrate` binary.

```sh
cd cmd/migrate
go build -v 
./migrate up
```

### HTTP server

Build and run the `httpserver` binary.

```sh
cd cmd/httpserver
go build -v
./httpserver
```
