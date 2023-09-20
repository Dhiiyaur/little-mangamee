package utils

import (
	"errors"
	"net/http"
)

type IdentityErr struct {
	HttpCode  int
	ErrorCode string
	Message   string
}

var (
	ERR_INTERNAL_SERVER = errors.New("internal server error")
	ERR_BAD_REQUEST     = errors.New("bad request")
	ERR_NOT_FOUND       = errors.New("not found")
)

var (
	ERR_MAP = map[error]IdentityErr{
		ERR_INTERNAL_SERVER: {
			HttpCode:  http.StatusInternalServerError,
			ErrorCode: "51",
			Message:   ERR_INTERNAL_SERVER.Error(),
		},
		ERR_BAD_REQUEST: {
			HttpCode:  http.StatusBadRequest,
			ErrorCode: "S01",
			Message:   ERR_BAD_REQUEST.Error(),
		},
		ERR_NOT_FOUND: {
			HttpCode:  http.StatusNotFound,
			ErrorCode: "S02",
			Message:   ERR_NOT_FOUND.Error(),
		},
	}
)

func FindError(rc string) error {

	var err_key error
	found := false

	for key, errData := range ERR_MAP {
		if errData.ErrorCode == rc {
			err_key = key
			found = true
			break
		}
	}
	if !found {
		err_key = errors.New("internal server error")
	}

	return err_key
}
