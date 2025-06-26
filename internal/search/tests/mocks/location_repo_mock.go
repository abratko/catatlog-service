package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/config/dc"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/location/app/dto"
)

type LocationRepoMock struct {
	mock.Mock
}

func (m *LocationRepoMock) FindByIds(ctx context.Context, ids []string) ([]dto.Location, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]dto.Location), args.Error(1)
}

func (m *LocationRepoMock) CancelReinit() {
	m.Called()
}

func NewLocationRepoMock() *LocationRepoMock {
	mockLocationRepo := &LocationRepoMock{}
	dc.LocationRepo.Mock(mockLocationRepo)

	return mockLocationRepo
}
