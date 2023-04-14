package gormsql

import (
	"errors"

	"gorm.io/gorm"
)

func Transaction() (*gorm.Tx, error) {
	return nil, errors.New("unimplemented")
}
