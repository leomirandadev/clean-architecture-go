
VERSION = $(shell git branch --show-current)
NAME = $(shell echo $(CURRENTNAME) | sed 's/\(.\)\([A-Z]\)/\1-\2/g' | tr '[:upper:]' '[:lower:]')

DB_HOST_MYSQL=127.0.0.1:3306
DB_USER_MYSQL=root
DB_PASSWORD_MYSQL=root
DB_NAME_MYSQL=golang_mysql

build: 
	./docker/build-image.sh $(NAME) $(VERSION) $(DB_HOST_MYSQL) $(DB_USER_MYSQL) $(DB_PASSWORD_MYSQL) $(DB_NAME_MYSQL)

run:
	DB_HOST_MYSQL=$(DB_HOST_MYSQL) DB_USER_MYSQL=$(DB_USER_MYSQL) DB_PASSWORD_MYSQL=$(DB_PASSWORD_MYSQL) DB_NAME_MYSQL=$(DB_NAME_MYSQL) VERSION=$(VERSION) go run main.go

run-watch:
	DB_HOST_MYSQL=$(DB_HOST_MYSQL) DB_USER_MYSQL=$(DB_USER_MYSQL) DB_PASSWORD_MYSQL=$(DB_PASSWORD_MYSQL) DB_NAME_MYSQL=$(DB_NAME_MYSQL) VERSION=$(VERSION) nodemon --exec go run main.go --signal SIGTERM

mig-create: 
	goose -dir ./migrations create $(MIG_NAME) sql 

mig-status: 
	goose mysql $(DB_CONNECTION) status

mig-up: 
	goose -dir ./migrations mysql $(DB_CONNECTION) up

mig-down: 
	goose -dir ./migrations mysql $(DB_CONNECTION) down
