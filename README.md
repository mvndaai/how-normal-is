# how-normal-is

App for hownormalis.com.

## Structure

- `/app`: Flutter client (web/iOS/android target)
- `/cmd/api` + `/internal`: Go backend services using `ctxerr`
- `/supabase/schema.sql`: Supabase/Postgres schema for auth-linked data

## Backend quick start

```bash
go test ./...
go run ./cmd/api
```
