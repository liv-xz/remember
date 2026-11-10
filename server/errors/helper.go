package errors

import (
	goerrors "errors"

	"gorm.io/gorm"
)

func IsRecordNotFound(err error) bool {
	return goerrors.Is(err, gorm.ErrRecordNotFound)
}
