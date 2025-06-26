package grpc

import (
	"encoding/json"
	"strconv"

	pb "gitlab.trgdev.com/gotrg/white-label/services/catalog/api/catalog/proto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
	"google.golang.org/protobuf/types/known/structpb"
)

func createFiltersResult(
	filters *model.Filters,
) *pb.FiltersResult {

	items := make([]*pb.Filter, len(filters.OnlyFaceted()))
	var i int
	for _, item := range filters.OnlyFaceted() {
		items[i] = createFilter(item)
		i++
	}

	return &pb.FiltersResult{
		Filters: items,
	}
}

func createFilter(filter *model.Filter) *pb.Filter {
	pbFilter := &pb.Filter{
		Key:   filter.Key(),
		Label: filter.Label(),
	}

	if filter.OperationType() == model.TermsOp {
		terms := filter.Terms()
		pbFilter.Terms = make([]*pb.TermOpt, len(terms))
		for i, term := range terms {
			pbFilter.Terms[i] = &pb.TermOpt{
				Value: term.Value,
				Count: uint32(term.Count),
			}

			if term.Props != nil {
				pbFilter.Terms[i].Props = createProps(term.Props)
			}
		}
	}

	if filter.OperationType() == model.RangeOp {
		pbFilter.Range = &pb.RangeOpt{
			From: strconv.FormatFloat(filter.Range().From, 'f', -1, 64),
			To:   strconv.FormatFloat(filter.Range().To, 'f', -1, 64),
		}
	}

	return pbFilter
}

func createProps(props map[string]any) *structpb.Struct {
	asJson, err := json.Marshal(props)
	if err != nil {
		return nil
	}

	var structSource map[string]any
	err = json.Unmarshal(asJson, &structSource)
	if err != nil {
		return nil
	}

	structValue, err := structpb.NewStruct(structSource)
	if err != nil {
		return nil
	}

	return structValue
}
