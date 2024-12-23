package middleware

import (
	"net/http"
	"testing"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/stretchr/testify/assert"
)

func TestEnsurelogId(t *testing.T) {
	t.Run("logId already exists", func(t *testing.T) {
		input := &openapi3filter.AuthenticationInput{
			RequestValidationInput: &openapi3filter.RequestValidationInput{
				Request: &http.Request{
					Header: http.Header{},
				},
			},
		}
		input.RequestValidationInput.Request.Header.Set(logHeaderString, "existing-log-id")

		logId := ensureLogID(input)

		assert.Equal(t, "existing-log-id", logId)
	})

	t.Run("logId generated if not exists", func(t *testing.T) {
		input := &openapi3filter.AuthenticationInput{
			RequestValidationInput: &openapi3filter.RequestValidationInput{
				Request: &http.Request{
					Header: http.Header{},
				},
			},
		}

		logId := ensureLogID(input)

		assert.NotEmpty(t, logId)
		assert.Equal(t, logId, input.RequestValidationInput.Request.Header.Get(logHeaderString))
	})
}
