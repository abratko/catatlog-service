package fields

import (
	"errors"

	"github.com/Jeffail/gabs/v2"
)

func CountItems(
	parsedJson *gabs.Container,
) (uint32, error) {

	count, ok := parsedJson.Path("doc_count").Data().(float64)
	if !ok {
		return 0, errors.New("casting value for CountItems to uint32 failed")
	}

	return uint32(count), nil
}
