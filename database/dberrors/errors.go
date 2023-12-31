package dberrors

import (
	"errors"
	"strings"

	"github.com/gogo/status"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

func ToStatusErr(err error, msgs ...string) error {
	if err == nil {
		return nil
	}

	msgs = append(msgs, "%v")
	message := strings.Join(msgs, ":")

	switch actualErr := err.(type) {
	case *pgconn.PgError:
		return fromPgConn(actualErr, message)

	default:

		code := codes.Internal

		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = codes.NotFound
		}

		return status.Errorf(code, message, err.Error())
	}
}

func fromPgConn(err *pgconn.PgError, message string) error {
	switch err.Code {
	case "23505":
		return status.Errorf(codes.AlreadyExists, message, err.Message)
	}

	return status.Errorf(codes.Unavailable, message, err.Message)
}
