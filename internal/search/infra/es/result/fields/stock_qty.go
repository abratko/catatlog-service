package fields

import (
	"github.com/Jeffail/gabs/v2"
)

func StockQty(
	parsedJson *gabs.Container,
) (uint32, error) {
	value, err := ValueAgg("facility.facility.stockQty", parsedJson)
	if err != nil {
		return 0, err
	}

	return uint32(value), nil
}
