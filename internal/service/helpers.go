package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/xyz-corp/xyz-football-api/pkg/apperror"
	"gorm.io/gorm"
)

func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// wrapNotFound converts a gorm.ErrRecordNotFound into a domain NotFound error,
// otherwise wraps as an internal error.
func wrapNotFound(err error, message string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.NotFound(message)
	}
	return apperror.Internal(err.Error())
}
