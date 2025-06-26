package app

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/location/app/dto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/location/app/errs"
)

type fetchCommand func(context.Context) (chan dto.Location, error)

type Repo struct {
	expAfter   time.Duration
	timer      *time.Timer
	mu         sync.RWMutex
	fetchError error
	fetch      fetchCommand
	byIds      map[string]*dto.Location
}

func NewRepo(
	fetch fetchCommand,
	expAfter string,
) *Repo {
	parsedDuration, err := time.ParseDuration(expAfter)
	if err != nil {
		panic(err)
	}

	repo := &Repo{
		expAfter: parsedDuration,
		fetch:    fetch,
	}

	repo.init()

	return repo
}

func (c *Repo) FindByIds(ctx context.Context, ids []string) ([]dto.Location, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.fetchError != nil {
		return nil, c.fetchError
	}

	result := make([]dto.Location, 0)

	if len(c.byIds) == 0 {
		return result, nil
	}

	for _, id := range ids {
		if location, ok := c.byIds[id]; ok {
			result = append(result, *location)
		}
	}

	return result, nil
}

func (c *Repo) CancelReinit() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.timer != nil {
		c.timer.Stop()
	}
}

func (c *Repo) init() {
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

		byIds := make(map[string]*dto.Location)

		for location := range generator {
			byIds[strconv.Itoa(location.Id)] = &location
		}

		c.timer = time.AfterFunc(c.expAfter, func() {
			c.init()
		})

		c.byIds = byIds
	}()

	wg.Wait()
}
