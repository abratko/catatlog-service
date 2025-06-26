package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/config/dc"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/category/app/dto"
)

type CategoryRepoMock struct {
	mock.Mock
}

func (m *CategoryRepoMock) FindByIds(ctx context.Context, ids []string) ([]dto.Category, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]dto.Category), args.Error(1)
}

func (m *CategoryRepoMock) CancelReinit() {
	m.Called()
}

func NewCategoryRepoMock() *CategoryRepoMock {
	mockCategoryRepo := &CategoryRepoMock{}
	dc.CategoryRepo.Mock(mockCategoryRepo)

	return mockCategoryRepo
}
