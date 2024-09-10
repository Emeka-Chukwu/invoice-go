# invoice-go
# Description
Building an invoicing spplication.


### Technical tools used

- Golang
- Docker
- Postgres (included in docker)
- Migrations
- jwt
- testify
- mock
- sqlmock

### Running

Clone the repo and run the below command at the project root 

```console
$ go mod tidy
```

Setting up all containers

```console
$ make up
```

### Create and drop DB


```console
$ make createdb

$ make dropdb
```
### Migration

Migrating sql file into db

```console
$ make migrateup
```

Dropping tables 

```console
$ make migratedown
```
Note: Run migration when docker is running and db migration is enabled to run at the start of the application

### Run Test

Run all the test in the project with

```console
$ make test
```