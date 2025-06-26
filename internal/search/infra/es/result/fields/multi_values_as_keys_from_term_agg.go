package fields

import (
	"fmt"

	"github.com/Jeffail/gabs/v2"
)

func MultiValuesAsKeysFromTermAgg(fieldName string, parsedJson *gabs.Container) ([]string, error) {
	buckets, ok := parsedJson.Path(fieldName + ".buckets").Data().([]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to cast buckets to array for %s", fieldName)
	}

	values := make([]string, len(buckets))
	for i, bucket := range buckets {
		bucket, ok := bucket.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to cast bucket to map for %s", fieldName)
		}
		value, ok := bucket["key"].(string)
		if !ok {
			return nil, fmt.Errorf("failed to cast value to string for %s", fieldName)
		}
		values[i] = value
	}

	return values, nil
}
