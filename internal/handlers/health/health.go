package health

import (
	"context"
	"net/http"

	"github.com/leomirandadev/clean-architecture-go/internal/handlers/middlewares"
	"github.com/leomirandadev/clean-architecture-go/pkg/healthcheck"
	"github.com/leomirandadev/clean-architecture-go/pkg/httprouter"
)

type Status struct {
	Health  bool `json:"health"`
	Details any  `json:"details"`
}

func Init(mid *middlewares.Middleware, router httprouter.Router, healthChecker healthcheck.Health) {
	ctr := controllers{
		healthChecker: healthChecker,
	}

	router.GET("/v1/health", mid.Auth.Public(ctr.Health))
}

type controllers struct {
	healthChecker healthcheck.Health
}

// health swagger document
// @Description Health checker
// @Tags health
// @Produce json
// @Success 200 {object} Status
// @Router /v1/health [get]
func (ctr controllers) Health(c httprouter.Context) error {
	return c.JSON(http.StatusOK, ctr.healthChecker.Health(context.Background()))
}
