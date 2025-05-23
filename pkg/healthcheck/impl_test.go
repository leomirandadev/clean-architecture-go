package healthcheck

import (
	"context"
	"errors"
	"slices"
	"sort"
	"testing"
)

func TestRegisterMany(t *testing.T) {
	testCases := map[string]struct {
		input  map[string]Checker
		output []HealthStatus
	}{
		"no checks": {
			input:  map[string]Checker{},
			output: []HealthStatus{},
		},
		"one check - all success": {
			input: map[string]Checker{
				"checker1": func(ctx context.Context) error {
					return nil
				},
			},
			output: []HealthStatus{
				{
					Key: "checker1",
					OK:  true,
				},
			},
		},
		"many checkers - all success": {
			input: map[string]Checker{
				"checker1": func(ctx context.Context) error {
					return nil
				},
				"checker2": func(ctx context.Context) error {
					return nil
				},
			},
			output: []HealthStatus{
				{
					Key: "checker1",
					OK:  true,
				},
				{
					Key: "checker2",
					OK:  true,
				},
			},
		},
		"many checkers - not all success": {
			input: map[string]Checker{
				"checker1": func(ctx context.Context) error {
					return nil
				},
				"checker2": func(ctx context.Context) error {
					return errors.New("connection error")
				},
			},
			output: []HealthStatus{
				{
					Key: "checker1",
					OK:  true,
				},
				{
					Key:    "checker2",
					OK:     false,
					ErrMsg: "connection error",
				},
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			h := New()
			for key, exec := range testCase.input {
				h = h.Register(key, exec)
			}

			output := h.Health(context.Background())

			// fix sorting
			sort.Slice(output, func(i, j int) bool {
				return output[i].Key < output[j].Key
			})

			if !slices.Equal(output, testCase.output) {
				t.Error("expects:", testCase.output, "given:", output)
			}
		})
	}
}
