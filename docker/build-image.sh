#!/bin/bash

NAME=$1
VERSION=$2
DB_HOST_MYSQL=$3
DB_USER_MYSQL=$4
DB_PASSWORD_MYSQL=$5
DB_NAME_MYSQL=$6

echo $NAME: Compilando o micro-serviço $NAME
# env GOOS=linux GOARCH=amd64 go build -o dist/$NAME
go build -o dist/$NAME

echo $NAME: Escrevendo o Dockerfile
CAT <<EOF > Dockerfile
    FROM alpine:3.7

    COPY ./dist/$NAME /opt/$NAME

    ENV DB_HOST_MYSQL=$DB_HOST_MYSQL DB_USER_MYSQL=$DB_USER_MYSQL DB_PASSWORD_MYSQL=$DB_PASSWORD_MYSQL DB_NAME_MYSQL=$DB_NAME_MYSQL VERSION=$VERSION

    WORKDIR /opt
    EXPOSE 5050

    ENTRYPOINT ./$NAME
EOF

echo $NAME: Construindo a imagem
docker build -t $NAME .

echo $NAME: Removendo artefatos
rm ./Dockerfile
rm -rf ./dist
