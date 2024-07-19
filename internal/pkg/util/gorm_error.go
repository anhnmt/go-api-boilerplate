package util

import (
	"errors"

	"github.com/bytedance/sonic"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var errCodes = map[string]error{
	"23505": gorm.ErrDuplicatedKey,
	"23503": gorm.ErrForeignKeyViolated,
	"42703": gorm.ErrInvalidField,
	"23514": gorm.ErrCheckConstraintViolated,
}

func GormTranslate(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if translatedErr, found := errCodes[pgErr.Code]; found {
			return translatedErr
		}
		return err
	}

	parsedErr, marshalErr := sonic.Marshal(err)
	if marshalErr != nil {
		return err
	}

	var errMsg postgres.ErrMessage
	unmarshalErr := sonic.Unmarshal(parsedErr, &errMsg)
	if unmarshalErr != nil {
		return err
	}

	if translatedErr, found := errCodes[errMsg.Code]; found {
		return translatedErr
	}
	return err
}
