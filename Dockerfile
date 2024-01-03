FROM alpine:3.7

WORKDIR /app

COPY ./dist/clean-architecture-go/api ./
COPY ./.env ./

# timezone
RUN apk update && apk --no-cache add tzdata 

EXPOSE 8080

ENTRYPOINT ./api
