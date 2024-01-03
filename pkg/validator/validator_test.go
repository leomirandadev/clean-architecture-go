package validator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	testCases := map[string]struct {
		input       any
		expectedErr error
	}{
		"Success: required": {
			input: struct {
				Info string `validate:"required"`
			}{
				Info: "result",
			},
			expectedErr: nil,
		},
		"Error: required": {
			input: struct {
				Info string `validate:"required"`
			}{},
			expectedErr: errors.New("Key: 'Info' Error:Field validation for 'Info' failed on the 'required' tag"),
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {

			err := Validate(testCase.input)
			if err == nil && testCase.expectedErr == err {
				return
			}

			if err != nil && testCase.expectedErr == nil || err == nil && testCase.expectedErr != nil {
				t.Errorf("actual:\n%v\nexpected:\n%v", err, testCase.expectedErr)
				return
			}

			assert.Equal(t, testCase.expectedErr.Error(), err.Error())
		})
	}
}
