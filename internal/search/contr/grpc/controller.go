package grpc

import (
	"context"
	"encoding/json"
	"errors"

	pb "gitlab.trgdev.com/gotrg/white-label/services/catalog/api/catalog/proto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/errs"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model/dto"
	coreGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fetchFilters = func(context.Context, map[string]any) (*model.Filters, error)
type fetchItems = func(context.Context, dto.SearchCriteria) (*dto.ItemsCollection, error)
type createCriteria = func(string) (*dto.SearchCriteria, error)

type Controller struct {
	pb.UnimplementedSearchServer
	// Use interface here because it is usefull fo tests
	// and we can inject a mock

	fetchItems          fetchItems
	fetchFilters        fetchFilters
	createItemsResult   func(*dto.ItemsCollection) *pb.ItemsResult
	createFiltersResult func(*model.Filters) *pb.FiltersResult
}

func NewController(
	fetchFilters fetchFilters,
	fetchItems fetchItems,
) *Controller {
	return &Controller{
		fetchItems:          fetchItems,
		fetchFilters:        fetchFilters,
		createItemsResult:   createItemsResult,
		createFiltersResult: createFiltersResult,
	}
}

func (c *Controller) Items(
	ctx context.Context,
	request *pb.ItemsRequest,
) (*pb.ItemsResult, error) {

	args := dto.SearchCriteria{}
	if request.Body != "" {
		err := json.Unmarshal([]byte(request.Body), &args)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid request body")
		}
	}

	appResult, err := c.fetchItems(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "app service error: %v", err)
	}

	return c.createItemsResult(appResult), nil
}

func (c *Controller) Filters(
	ctx context.Context,
	request *pb.FiltersRequest,
) (*pb.FiltersResult, error) {

	args := map[string]any{}
	if request.Body != "" {
		err := json.Unmarshal([]byte(request.Body), &args)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid request body")
		}
	}
	filters, err := c.fetchFilters(ctx, args)
	if err != nil {
		return nil, c.createStatusError(err)
	}

	return c.createFiltersResult(filters), nil
}

func (s *Controller) Start(server coreGrpc.ServiceRegistrar) {
	pb.RegisterSearchServer(server, s)
}

func (c *Controller) createStatusError(appErr error) error {
	if errors.Is(appErr, errs.ErrInternal) {
		return status.Errorf(codes.Internal, appErr.Error())
	}

	if errors.Is(appErr, errs.ErrInvalidArg) {
		return status.Errorf(codes.InvalidArgument, appErr.Error())
	}

	return status.Errorf(codes.Internal, appErr.Error())
}
