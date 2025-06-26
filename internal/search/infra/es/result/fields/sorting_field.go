package fields

import (
	"fmt"

	"github.com/Jeffail/gabs/v2"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model/dto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/infra/es/consts"
)

func SortingField(
	parsedJson *gabs.Container,
) (*dto.SortingField, error) {

	var sortingField *dto.SortingField

	for _, fieldName := range consts.SortingFields {
		if !parsedJson.ExistsP(fieldName + ".value") {
			continue
		}

		sortingField = &dto.SortingField{
			Field: fieldName,
		}

		anyValue := parsedJson.Path(fieldName + ".value").Data()
		if anyValue == nil {
			break
		}

		sortingField.Value = fmt.Sprintf("%v", anyValue)

		break
	}

	return sortingField, nil
}
