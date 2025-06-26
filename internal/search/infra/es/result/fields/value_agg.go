package fields

import (
	"fmt"

	"github.com/Jeffail/gabs/v2"
)

func ValueAgg(fieldName string, parsedJson *gabs.Container) (float64, error) {
	value, ok := parsedJson.Path(fieldName + ".value").Data().(float64)
	if !ok {
		return 0, fmt.Errorf("failed to cast value to float64 for %s", fieldName)
	}

	return value, nil
}
