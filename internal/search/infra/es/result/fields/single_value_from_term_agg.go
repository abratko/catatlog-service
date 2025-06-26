package fields

import (
	"fmt"

	"github.com/Jeffail/gabs/v2"
)

func SingleValueFromTermAgg(fieldName string, parsedJson *gabs.Container) (string, error) {

	value, ok := parsedJson.Path(fieldName + ".buckets.0.key").Data().(string)
	if !ok {
		return "", fmt.Errorf("failed to cast value to string for %s", fieldName)
	}

	return value, nil
}
