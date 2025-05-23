# build step
FROM golang:1.23-alpine AS base
WORKDIR /app
COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o api cmd/api/*.go 

# running step
FROM alpine:3.21.3 AS production
WORKDIR /app
COPY --from=base /app/api .
COPY --from=base /app/.env .
RUN apk update && apk --no-cache add tzdata 
EXPOSE 8080
ENTRYPOINT ["./api"]
