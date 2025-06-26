package fields

import (
	"fmt"

	"github.com/Jeffail/gabs/v2"
)

func Facilities(
	parsedJson *gabs.Container,
) (map[string]uint32, error) {
	if !parsedJson.ExistsP("facility.facility.facilityAddressName.buckets") {
		return nil, fmt.Errorf("not found facility buckets in es response")
	}

	buckets := parsedJson.Path("facility.facility.facilityAddressName.buckets").Children()
	if len(buckets) == 0 {
		return nil, fmt.Errorf("failed to cast buckets to array for facility")
	}

	values := make(map[string]uint32)

	for _, item := range buckets {
		value, ok := item.Path("key").Data().(string)
		if !ok {
			return nil, fmt.Errorf("failed to cast value to string for facility")
		}
		count, err := ValueAgg("stockQty", item)
		if err != nil {
			return nil, err
		}
		values[value] = uint32(count)
	}

	return values, nil
}
