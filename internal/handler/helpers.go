package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xyz-corp/xyz-football-api/pkg/apperror"
	"github.com/xyz-corp/xyz-football-api/pkg/response"
	"github.com/xyz-corp/xyz-football-api/pkg/validator"
)

func bindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		response.Error(c, apperror.Validation("Validasi gagal", validator.Details(err)...))
		return false
	}
	return true
}

func paramUUID(c *gin.Context, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		response.Error(c, apperror.Validation("Format UUID tidak valid"))
		return uuid.Nil, false
	}
	return id, true
}

// handleErr writes the appropriate error envelope, converting unexpected errors to 500.
func handleErr(c *gin.Context, err error) {
	if appErr, ok := err.(*apperror.AppError); ok {
		response.Error(c, appErr)
		return
	}
	response.Error(c, apperror.Internal(err.Error()))
}
