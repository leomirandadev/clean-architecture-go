# 🚀 Getting Started

## ⭐ Dependencies
You must need to install some dependencies before start this project. These dependencies are:
- Go version 1.21 or greater; [🔗 https://go.dev/dl/](https://go.dev/dl/)
- Docker and Docker Compose; [🔗 https://docs.docker.com/compose/install/](https://docs.docker.com/compose/install/)

## ⭐ Installing
To install the golang dependencies you must run:
```shell
make install
```

## ⭐ Environment variables
To set the environment variables you only need to create a file ``.env`` based on ``.env.example`` and update the variables when necessary.

## ⭐ Running
To have ``database``, ``cache`` and ``jaeger`` on your local machine you only need to run:
```shell
make local-up
```

To run the migrations you can run:
```shell
make mig-up
```

and to get the api running you can run:
```shell
make run
# or
make run-watch
```
*Ps.: run-watch will watch the files changes and restart the API when it was necessary.*

## ⭐ Endpoint docs (Swagger)
To access the swagger you can just run:
```shell
make open-swagger
```
or just access:  [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

When you create a new endpoint, we need to run this command to update the swagger docs:
```shell
make docs
```

[🔗 Here](https://github.com/swaggo/swag) you can see more details about the library this project is using to generate this docs based on comments.

## ⭐ Tracing (Jaeger)
To access the jaeger you can just run:
```shell
make open-jaeger
```
or just access:  [http://localhost:16686/search](http://localhost:16686/search)

## ⭐ Migrations

Migrations is responsible for versionate the database and you can see all migrations that we had until now on the [migrations](/migrations) folder.

### Running all migrations
To run the migrations you can run:
```shell
make mig-up
```
This will run the migrations that your database don't have yet.

### Rollback
If something goes wrong and you need to make a rollback on the last migration, you can run:
```shell
make mig-down
```

### Creating a new migration
Do you want to create a new database change? So you only need to run:
```shell
make mig-create MIG_NAME=migration_name
```
where you must update the ``migration_name`` with something that represents your change. One file is going to be created in the [migrations](/migrations) folder and you'll need to put the query that you want and the rollback query.

To know more about how to write your migration, you can go to [goose docs](https://github.com/pressly/goose).


## ⭐ Test and Mock

### Running the test pipeline
To run the pipeline test we have two commands:

```shell
make test
```
This command runs the whole pipeline.

The second one runs the pipeline too, but it also shows on a browser the details and where is the uncoverage code is. The command to do this is:

```shell
make test-cover
```
### Mocking a package
To mock a package or file this project is using the ``go:generate`` + ``gomock`` and you can do it only including in the first line of the go file the comment:

```go
//go:generate mockgen -source file_name.go -destination mocks/cache_mock.go -package mocks
```
*Ps.: Don't forget to update file_name.go to the right filename.*

After that, you only need to run:
```shell
make mock
```
and if you invert the dependency currectly [Gomock](https://github.com/golang/mock) will create a folder with the mocks that you need.
