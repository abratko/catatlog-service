package errs

import (
	"errors"
)

var ErrFetchItems = errors.New("failed to fetch items")
var ErrItemsNotFound = errors.New("not found")
