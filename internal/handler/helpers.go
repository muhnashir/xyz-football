package handler

import (
	"strings"

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

// publicBaseURL derives the scheme+host the current request actually came in on,
// so links we hand back always match the port/host the client used — instead of
// relying on a separately configured (and easily out-of-sync) base URL.
func publicBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	host := c.Request.Host
	if fwd := c.GetHeader("X-Forwarded-Host"); fwd != "" {
		host = fwd
	}
	return scheme + "://" + host
}

// resolveLogoURL turns a stored relative path (e.g. "teams/logo/x.jpg") into an
// absolute URL against the current request's host. Values that are already
// absolute (legacy rows, or a client-supplied external URL) are left untouched.
func resolveLogoURL(c *gin.Context, logoURL *string) *string {
	if logoURL == nil || *logoURL == "" {
		return logoURL
	}
	if strings.HasPrefix(*logoURL, "http://") || strings.HasPrefix(*logoURL, "https://") {
		return logoURL
	}
	full := publicBaseURL(c) + "/" + strings.TrimPrefix(*logoURL, "/")
	return &full
}
