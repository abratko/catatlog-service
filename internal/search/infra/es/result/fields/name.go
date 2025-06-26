package fields

import (
	"fmt"

	"github.com/Jeffail/gabs/v2"
)

func Name(
	parsedJson *gabs.Container,
) (string, error) {

	name, ok := parsedJson.Path("key").Data().(string)
	if !ok {
		return "", fmt.Errorf("failed to cast value to string for name")
	}

	return name, nil
}
