package rest_err

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestErr_Error(t *testing.T) {
	errMessage := "test error message"
	restErr := &RestErr{Message: errMessage}
	assert.Equal(t, errMessage, restErr.Error(), "Error() should return the expected message")
}

func TestNewBadRequestError(t *testing.T) {
	message := "test bad request message"
	restErr := NewBadRequestError(message)

	assert.NotNil(t, restErr, "NewBadRequestError should return a non-nil error")
	assert.Equal(t, message, restErr.Message, "Message should be set correctly")
	assert.Equal(t, "bad_request", restErr.Err, "Err should be set to 'bad_request'")
	assert.Equal(t, http.StatusBadRequest, restErr.Code, "Code should be set to http.StatusBadRequest")
	assert.Empty(t, restErr.Causes, "Causes should be an empty slice")
}

func TestNewUnauthorizedRequestError(t *testing.T) {
	message := "test unauthorized message"
	restErr := NewUnauthorizedRequestError(message)

	assert.NotNil(t, restErr, "NewUnauthorizedRequestError should return a non-nil error")
	assert.Equal(t, message, restErr.Message, "Message should be set correctly")
	assert.Equal(t, "unauthorized", restErr.Err, "Err should be set to 'unauthorized'")
	assert.Equal(t, http.StatusUnauthorized, restErr.Code, "Code should be set to http.StatusUnauthorized")
	assert.Empty(t, restErr.Causes, "Causes should be an empty slice")
}

func TestNewBadRequestValidationError(t *testing.T) {
	message := "test validation error message"
	causes := []Causes{{Field: "field1", Message: "validation error"}}
	restErr := NewBadRequestValidationError(message, causes)

	assert.NotNil(t, restErr, "NewBadRequestValidationError should return a non-nil error")
	assert.Equal(t, message, restErr.Message, "Message should be set correctly")
	assert.Equal(t, "bad_request", restErr.Err, "Err should be set to 'bad_request'")
	assert.Equal(t, http.StatusBadRequest, restErr.Code, "Code should be set to http.StatusBadRequest")
	assert.Equal(t, causes, restErr.Causes, "Causes should be set correctly")
}

func TestNewInternalServerError(t *testing.T) {
	message := "test internal server error message"
	restErr := NewInternalServerError(message)

	assert.NotNil(t, restErr, "NewInternalServerError should return a non-nil error")
	assert.Equal(t, message, restErr.Message, "Message should be set correctly")
	assert.Equal(t, "internal_server_error", restErr.Err, "Err should be set to 'internal_server_error'")
	assert.Equal(t, http.StatusInternalServerError, restErr.Code, "Code should be set to http.StatusInternalServerError")
	assert.Empty(t, restErr.Causes, "Causes should be an empty slice")
}

func TestNewNotFoundError(t *testing.T) {
	message := "test not found error message"
	restErr := NewNotFoundError(message)

	assert.NotNil(t, restErr, "NewNotFoundError should return a non-nil error")
	assert.Equal(t, message, restErr.Message, "Message should be set correctly")
	assert.Equal(t, "not_found", restErr.Err, "Err should be set to 'not_found'")
	assert.Equal(t, http.StatusNotFound, restErr.Code, "Code should be set to http.StatusNotFound")
	assert.Empty(t, restErr.Causes, "Causes should be an empty slice")
}

func TestNewForbiddenError(t *testing.T) {
	message := "test forbidden error message"
	restErr := NewForbiddenError(message)

	assert.NotNil(t, restErr, "NewForbiddenError should return a non-nil error")
	assert.Equal(t, message, restErr.Message, "Message should be set correctly")
	assert.Equal(t, "forbidden", restErr.Err, "Err should be set to 'forbidden'")
	assert.Equal(t, http.StatusForbidden, restErr.Code, "Code should be set to http.StatusForbidden")
	assert.Empty(t, restErr.Causes, "Causes should be an empty slice")
}
