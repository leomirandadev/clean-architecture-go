package aescrypt

import (
	"bytes"
	"errors"
	"testing"
)

func Test1(t *testing.T) {
	testCases := map[string]struct {
		inputKey     string
		inputMessage []byte
		outputErr    error
	}{
		"success case - no errors": {
			inputKey:     "1f345f16b10298b7cf3ce8411cae52d03533de4f47516cb3673c530abbb83aad",
			inputMessage: []byte("this is the message"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			app := New(testCase.inputKey)

			result, err := app.Encrypt(testCase.inputMessage)
			if !errors.Is(testCase.outputErr, err) {
				t.Error("Errors doesn't match 1:", err, testCase.outputErr)
				return
			}

			message, err := app.Decrypt(result)
			if !errors.Is(testCase.outputErr, err) {
				t.Error("Errors doesn't match 2:", err, testCase.outputErr)
				return
			}

			if !bytes.Equal(message, testCase.inputMessage) {
				t.Error("result not equal", message, testCase.inputMessage)
			}
		})
	}
}
