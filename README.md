# Clean Architecture in Golang
The main goal of this project is create a structure to reuse that implement clean Architecture principals, allowed you change all the libraries and mock all de layers whenever you want.

## What technologies were implemented 
- Cache: Redis and Memcached;
- Hashing: Bcrypt;
- Http Router: Mux;
- Logger: Logrus;
- Mailer: net/stmp (native of go);
- token generator: JWT GO;
- DataBase manipulator: Sqlx and Gorm;

## Next Implementations
- Tests: [Gomock](https://github.com/golang/mock);
- Router: [Gin](https://github.com/gin-gonic/gin);
- Remote configuration: [Viper](https://github.com/spf13/viper);
- Message Broker: [Nats.io](https://github.com/nats-io/nats.go) and [Kafka](https://github.com/confluentinc/confluent-kafka-go);
- RPC layer: [GRPC](https://pkg.go.dev/google.golang.org/grpc);
- DataBase manipulator: [mongodb](https://github.com/mongodb/mongo-go-driver);
- Observability: [openTelemetry](https://github.com/open-telemetry/opentelemetry-go);

## Minimum softwares
- [Docker](https://docs.docker.com/desktop/);
- [Golang](https://golang.org/doc/install);
- [Nodemon](https://nodemon.io/);
- [goose](https://github.com/pressly/goose);