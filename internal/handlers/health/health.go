package health

import (
	"net/http"

	"github.com/leomirandadev/clean-architecture-go/internal/handlers/middlewares"
	"github.com/leomirandadev/clean-architecture-go/pkg/httprouter"
)

type Status struct {
	Health  bool `json:"health"`
	Details any  `json:"details"`
}

func Init(mid *middlewares.Middleware, router httprouter.Router) {
	ctr := controllers{}
	router.GET("/v1/health", mid.Auth.Public(ctr.Health))
}

type controllers struct{}

// health swagger document
// @Description Health checker
// @Tags health
// @Produce json
// @Success 200 {object} Status
// @Router /v1/health [get]
func (ctr controllers) Health(c httprouter.Context) error {
	// TODO implementing the checkers
	return c.JSON(http.StatusOK, Status{
		Health: true,
	})
}
