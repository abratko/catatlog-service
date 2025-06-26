package grpc

import (
	pb "gitlab.trgdev.com/gotrg/white-label/services/catalog/api/catalog/proto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model/dto"
)

func createItemsResult(
	appResult *dto.ItemsCollection,
) *pb.ItemsResult {

	items := make([]*pb.Item, len(appResult.Items()))
	for i, item := range appResult.Items() {
		items[i] = createItem(item)
	}

	return &pb.ItemsResult{
		TotalCount: uint32(appResult.TotalCount()),
		Items:      items,
	}
}

func createItem(appItem *dto.Item) *pb.Item {
	item := &pb.Item{
		Id:                appItem.Id,
		Name:              appItem.Name,
		Slug:              appItem.Slug,
		Image:             appItem.Image,
		LowestAsk:         float64(appItem.LowestAsk),
		StockQty:          uint32(appItem.StockQty),
		Conditions:        appItem.Conditions,
		Facilities:        appItem.Facilities,
		ProductDisclaimer: appItem.ProductDisclaimer,
		ProductSlug:       appItem.ProductSlug,
		ProductName:       appItem.ProductName,
		Skus:              appItem.Skus,
	}

	if appItem.SortingField != nil {
		item.SortingField = &pb.SortingField{
			Field: appItem.SortingField.Field,
			Value: appItem.SortingField.Value,
		}
	}

	return item
}
