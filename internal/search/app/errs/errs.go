package errs

import (
	"errors"
)

var ErrInvalidArg = errors.New("invalid argument")
var ErrInternal = errors.New("internal error")
