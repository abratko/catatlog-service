package app

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/category/app/dto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/category/app/errs"
)

type fetchCategories = func(context.Context) (chan dto.Category, error)

type CategoryRepo struct {
	expAfter time.Duration
	timer    *time.Timer
	mu       sync.RWMutex

	fetchError error
	fetch      fetchCategories

	asTree map[string]*dto.Category
	byIds  map[string]*dto.Category
	byPath map[string]*dto.Category
}

func NewCategoryRepo(
	fetch fetchCategories,
	expAfter string,
) *CategoryRepo {

	parsedDuration, err := time.ParseDuration(expAfter)
	if err != nil {
		panic(err)
	}

	repo := &CategoryRepo{
		expAfter: parsedDuration,
		fetch:    fetch,
	}

	repo.init()

	return repo
}

func (c *CategoryRepo) FindByIds(ctx context.Context, ids []string) ([]dto.Category, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.fetchError != nil {
		return nil, c.fetchError
	}

	result := make([]dto.Category, 0)

	if len(c.byIds) == 0 {
		return result, nil
	}

	for _, id := range ids {
		if category, ok := c.byIds[id]; ok {
			result = append(result, *category)
		}
	}

	return result, nil
}

func (c *CategoryRepo) CancelReinit() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.timer != nil {
		c.timer.Stop()
	}
}

func (c *CategoryRepo) init() {
	c.mu.Lock()
	defer c.mu.Unlock()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		generator, err := c.fetch(context.Background())
		if err != nil {
			c.fetchError = fmt.Errorf("%w: %s", errs.ErrFetchItems, err)
			return
		}

		byIds := make(map[string]*dto.Category)
		byPath := make(map[string]*dto.Category)
		asTree := make(map[string]*dto.Category)

		for category := range generator {
			byIds[category.Id] = &category
			byPath[strings.Join(category.Path, "/")] = &category
			asTree[category.Id] = &category
			// generator renders categories in order of their hierarchy: top-level categories first
			if parent := byIds[category.ParentId]; parent != nil {
				parent.Children = append(parent.Children, &category)
			}
		}

		c.timer = time.AfterFunc(c.expAfter, func() {
			c.init()
		})

		c.byIds = byIds
		c.byPath = byPath
		c.asTree = asTree
	}()

	wg.Wait()
}
