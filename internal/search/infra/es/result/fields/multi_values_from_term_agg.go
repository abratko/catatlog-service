package fields

import (
	"fmt"

	"github.com/Jeffail/gabs/v2"
)

func MultiValuesFromTermAgg(fieldName string, parsedJson *gabs.Container) (map[string]uint32, error) {
	buckets, ok := parsedJson.Path(fieldName + ".buckets").Data().([]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to cast buckets to array for %s", fieldName)
	}

	values := make(map[string]uint32)

	for _, item := range buckets {
		bucket, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to cast bucket to map for %s", fieldName)
		}

		value, ok := bucket["key"].(string)
		if !ok {
			return nil, fmt.Errorf("failed to cast value to string for %s", fieldName)
		}
		count, ok := bucket["doc_count"].(float64)
		if !ok {
			return nil, fmt.Errorf("failed to cast count to int for %s", fieldName)
		}
		values[value] = uint32(count)
	}

	return values, nil
}
