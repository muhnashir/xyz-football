package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyz-corp/xyz-football-api/pkg/apperror"
)

type Envelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type ErrorBody struct {
	Code    string            `json:"code"`
	Details []apperror.Detail `json:"details,omitempty"`
}

type ErrorEnvelope struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Error   ErrorBody `json:"error"`
}

func Success(c *gin.Context, httpStatus int, message string, data interface{}, meta interface{}) {
	c.JSON(httpStatus, Envelope{Success: true, Message: message, Data: data, Meta: meta})
}

func OK(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusOK, message, data, nil)
}

func Created(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusCreated, message, data, nil)
}

func Paginated(c *gin.Context, message string, data interface{}, meta interface{}) {
	Success(c, http.StatusOK, message, data, meta)
}

func Error(c *gin.Context, err *apperror.AppError) {
	c.JSON(err.HTTPStatus, ErrorEnvelope{
		Success: false,
		Message: err.Message,
		Error:   ErrorBody{Code: err.Code, Details: err.Details},
	})
}
