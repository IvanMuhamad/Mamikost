package models

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrAccessForbidden     = errors.New("access forbidden")
	ErrUserAlreadyExist    = errors.New("username or phone already exist")
	ErrDataNotFound        = errors.New("data not found")
	ErrInvalidUserPassword = errors.New("invalid username or password")
	ErrUserPasswordEmpty   = errors.New("username or password must not be empty")
	ErrLoginFailed         = errors.New("login failed")
	ErrInvalidToken        = errors.New("invalid token")
	ErrUnauthorizedToken   = errors.New("unAuthorized token")
	ErrFailedToAuthorized  = errors.New("UnAuthorized token")
	ErrFailedGenerateToken = errors.New("failed to generate token")
	ErrUpdateToken         = errors.New("update token failed")
	ErrNoRows              = errors.New("no rows in result set")
	/*category*/
	ErrCategoryNameDuplicate = errors.New("duplicate category name")
	ErrCategoryNotFound      = errors.New("category not found")
)

type Error struct {
	Errors map[string]interface{} `json:"errors"`
}

func NewValidationError(err error) *Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = err.Error()
	return &e
}

func NewError(err error) *Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["message"] = err.Error()
	return &e
}

func Nullable[T any](row *T, err error) (*T, error) {
	if err == nil {
		return row, nil
	}

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	return nil, err
}

func NullableList[T any](rows []*T, err error) ([]*T, error) {
	if err == nil {
		return rows, nil
	}

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	return nil, err
}

func NullableID(row string, err error) (string, error) {
	if err == nil {
		return row, nil
	}

	if err == pgx.ErrNoRows {
		return "", nil
	}

	return "", err
}
