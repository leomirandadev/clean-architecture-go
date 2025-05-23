package healthcheck

import "context"

type Health interface {
	Health(ctx context.Context) []HealthStatus
	Register(key string, exec Checker) Health
}

type HealthStatus struct {
	Key    string `json:"key"`
	OK     bool   `json:"ok"`
	ErrMsg string `json:"error_msg,omitempty"`
}

type Checker = func(ctx context.Context) error

type HealthChecker struct {
	Key     string
	Checker Checker
}
