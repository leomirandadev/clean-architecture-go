package swagger

import (
	docs "github.com/leomirandadev/clean-architecture-go/docs"
	"github.com/leomirandadev/clean-architecture-go/internal/handlers/middlewares"
	"github.com/leomirandadev/clean-architecture-go/pkg/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Init(mid *middlewares.Middleware, router httprouter.Router) {

	docs.SwaggerInfo.Title = "Swagger"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router.GET("/swagger/*", mid.Auth.BasicAuth(
		router.ParseHandler(
			httpSwagger.Handler(
				httpSwagger.URL("doc.json"),
			),
		),
	))
}
