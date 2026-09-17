package service

import "github.com/xyz-corp/xyz-football-api/pkg/apperror"

func apperrorDetail(field, message string) apperror.Detail {
	return apperror.Detail{Field: field, Message: message}
}
