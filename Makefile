
VERSION = $(shell git branch --show-current)
NAME = $(shell echo $(CURRENTNAME) | sed 's/\(.\)\([A-Z]\)/\1-\2/g' | tr '[:upper:]' '[:lower:]')

CACHE_URL=127.0.0.1:11211
CACHE_EXP=0
DB_CONNECTION = "root:root@(127.0.0.1:3306)/golang_mysql?charset=utf8\u0026readTimeout=30s\u0026writeTimeout=30s&parseTime=true&loc=Local"

build: 
	./docker/build-image.sh $(NAME) $(VERSION) $(DB_CONNECTION) $(CACHE_URL) $(CACHE_EXP)

run:
	DB_CONNECTION=$(DB_CONNECTION) CACHE_URL=$(CACHE_URL) CACHE_EXP=$(CACHE_EXP) VERSION=$(VERSION) go run main.go

run-watch:
	DB_CONNECTION=$(DB_CONNECTION) CACHE_URL=$(CACHE_URL) CACHE_EXP=$(CACHE_EXP) VERSION=$(VERSION) nodemon --exec go run main.go --signal SIGTERM

mig-create: 
	goose -dir ./migrations create $(MIG_NAME) sql 

mig-status: 
	goose mysql $(DB_CONNECTION) status

mig-up: 
	goose -dir ./migrations mysql $(DB_CONNECTION) up

mig-down: 
	goose -dir ./migrations mysql $(DB_CONNECTION) down
