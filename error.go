package chat

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrCode int

const (
	ERR_INTERNAL  = 0 // return 500
	ERR_DUPLICATE = 1 // return 422
	ERR_AUTHORIZE = 2 // return 401
	ERR_INVALID   = 3 // return 400
)

var StatusCodes = map[int]int{
	ERR_INTERNAL:  http.StatusInternalServerError,
	ERR_DUPLICATE: http.StatusUnprocessableEntity,
	ERR_AUTHORIZE: http.StatusUnauthorized,
	ERR_INVALID:   http.StatusBadRequest,
}

type Error struct {
	code int
	msg  string
}

func NewError(code int, msg string) *Error {
	return &Error{code: code, msg: msg}
}

func (e Error) Error() string { return e.msg }

func ResponseError(g *gin.Context, e error) {
	var err Error
	if errors.As(e, &err) {
		g.JSON(StatusCodes[err.code], err.msg)
	}
	g.JSON(http.StatusInternalServerError, e.Error())
}
