package healthcheck

import (
	"context"
)

func New() Health {
	return &impl{healthCheckers: make([]HealthChecker, 0)}
}

type impl struct {
	healthCheckers []HealthChecker
}

func (i impl) Health(ctx context.Context) []HealthStatus {
	result := make([]HealthStatus, 0, len(i.healthCheckers))

	for _, v := range i.healthCheckers {
		check := HealthStatus{
			Key: v.Key,
			OK:  true,
		}

		if err := v.Checker(ctx); err != nil {
			check.ErrMsg = err.Error()
			check.OK = false
		}

		result = append(result, check)
	}

	return result
}

func (i *impl) Register(key string, exec Checker) Health {
	i.healthCheckers = append(i.healthCheckers, HealthChecker{
		Key:     key,
		Checker: exec,
	})

	return i
}
