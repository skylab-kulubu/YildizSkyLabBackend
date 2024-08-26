# Yildiz Sky Lab

It is the backend project of the Yildiz Sky Lab website.

## Installations

Database migration tool [migrate](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md)

Sql queryies tool [sqlc](https://sqlc.dev)


## Migration Commands
  ### Migrate Up
  ```bash
migrate -path src/db/migration -database "postgresql url" -verbose up
```
  ### Migrate Down
```bash
migrate -path src/db/migration -database "postgresql url" -verbose down
```

## Update SQL Queries
  ### Generate
  ```bash
  sqlc generate
  ```

## Run Application
  ### build
  ```bash
go build -o bin/fs
  ```
  ### run
  ```bash
./bin/fs
