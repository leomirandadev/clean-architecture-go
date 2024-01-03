# build step
FROM golang:1.21 AS base
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api cmd/api/*.go 

# running step
FROM alpine:3.7 AS production
WORKDIR /app
COPY --from=base /app/api .
COPY --from=base /app/.env .
RUN apk update && apk --no-cache add tzdata 
EXPOSE 8080
ENTRYPOINT ./api
