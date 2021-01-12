
VERSION = $(shell git branch --show-current)
DB_CONNECTION = "root:root@(127.0.0.1:3306)/golang_mysql?charset=utf8\u0026readTimeout=30s\u0026writeTimeout=30s&parseTime=true&loc=Local"
DB_ENVS = DB_HOST_MYSQL=127.0.0.1:3306 DB_USER_MYSQL=root DB_PASSWORD_MYSQL=root DB_NAME_MYSQL=golang_mysql

run:
	$(DB_ENVS) VERSION=$(VERSION) go run main.go

run-watch:
	$(DB_ENVS) VERSION=$(VERSION) nodemon --exec go run main.go --signal SIGTERM

mig-create: 
	goose -dir ./migrations create $(MIG_NAME) sql 

mig-status: 
	goose mysql $(DB_CONNECTION) status

mig-up: 
	goose -dir ./migrations mysql $(DB_CONNECTION) up

mig-down: 
	goose -dir ./migrations mysql $(DB_CONNECTION) down
