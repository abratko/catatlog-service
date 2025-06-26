package result

import (
	"bufio"
	"fmt"

	"github.com/Jeffail/gabs/v2"
	"github.com/elastic/go-elasticsearch/v7/esapi"

	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model/dto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/infra/es/result/fields"
)

func Create(
	esResult *esapi.Response,
) (collection *dto.ItemsCollection, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("invalid response format from es: [%w]", err)
		}
	}()

	reader := bufio.NewReader(esResult.Body)

	parsedJSON, err := gabs.ParseJSONBuffer(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json: [%w]", err)
	}

	totalCount, ok := parsedJSON.Path("aggregations.count.value").Data().(float64)
	if !ok {
		return nil, fmt.Errorf("invalid format for totalCount")
	}

	buckets := parsedJSON.Path("aggregations.name.buckets").Children()
	items := make([]*dto.Item, len(buckets))

	for i, bucket := range buckets {
		items[i], err = createFamilyTail(bucket)
		if err != nil {
			return nil, err
		}
	}

	return dto.NewFamilyTailCollection(
		uint32(totalCount),
		items,
	), nil
}

func createFamilyTail(item *gabs.Container) (*dto.Item, error) {

	familyTail := dto.Item{}

	var err error
	familyTail.Id, err = fields.SingleValueFromTermAgg("id", item)
	if err != nil {
		return nil, err
	}

	familyTail.Name, err = fields.Name(item)
	if err != nil {
		return nil, err
	}

	familyTail.Skus, err = fields.MultiValuesAsKeysFromTermAgg("skus", item)
	if err != nil {
		return nil, err
	}

	familyTail.Slug, err = fields.SingleValueFromTermAgg("slug", item)
	if err != nil {
		return nil, err
	}

	familyTail.Image, err = fields.SingleValueFromTermAgg("image", item)
	if err != nil {
		return nil, err
	}

	familyTail.LowestAsk, err = fields.ValueAgg("lowestAsk", item)
	if err != nil {
		return nil, err
	}

	stockQty, err := fields.ValueAgg("facility.facility.stockQty", item)
	if err != nil {
		return nil, err
	}
	familyTail.StockQty = uint32(stockQty)

	familyTail.Conditions, err = fields.MultiValuesFromTermAgg("condition", item)
	if err != nil {
		return nil, err
	}

	familyTail.Facilities, err = fields.Facilities(item)
	if err != nil {
		return nil, err
	}

	familyTail.ProductDisclaimer, err = fields.SingleValueFromTermAgg("productDisclaimer", item)
	if err != nil {
		return nil, err
	}

	familyTail.ProductSlug, err = fields.SingleValueFromTermAgg("productSlug", item)
	if err != nil {
		return nil, err
	}

	familyTail.ProductName, err = fields.SingleValueFromTermAgg("productName", item)
	if err != nil {
		return nil, err
	}

	familyTail.SortingField, err = fields.SortingField(item)
	if err != nil {
		return nil, err
	}

	return &familyTail, nil
}
