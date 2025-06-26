package es

import (
	"bufio"
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/category/app/dto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/category/infra/es/request"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/pkg/es"
)

type FetchCategories struct {
	db *es.ReadAdapter
}

func NewFetchCategories(
	db *es.ReadAdapter,
) *FetchCategories {
	return &FetchCategories{
		db: db,
	}
}
func (c *FetchCategories) Exec(
	ctx context.Context,
) (chan dto.Category, error) {
	requestDto := request.CreateFetchCategories()

	res, err := c.db.Read(ctx, requestDto)
	if err != nil {
		return nil, err
	}

	return c.CreateCategoriesGenerator(res)
}

func (c *FetchCategories) CreateCategoriesGenerator(res *esapi.Response) (chan dto.Category, error) {
	reader := bufio.NewReader(res.Body)

	parsedJSON, err := gabs.ParseJSONBuffer(reader)
	if err != nil {
		return nil, errors.New("invalid response format from es")
	}

	hits := parsedJSON.Path("hits.hits").Children()
	if len(hits) == 0 {
		return nil, errors.New("not found hits in es response")
	}

	stack := []*gabs.Container{hits[0]} // Стек для хранения узлов

	ch := make(chan dto.Category)
	go func() {
		// Закрываем канал после завершения отправки данных
		defer close(ch)

		for len(stack) > 0 {
			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1] // Извлекаем последний элемент из стека

			if current.ExistsP("_source") {
				current = current.Path("_source")
			}

			category := c.createCategory(current)

			if current.ExistsP("children_data") {
				stack = append(stack, current.Path("children_data").Children()...)
			}

			ch <- category
		}
	}()

	return ch, nil
}
func (c *FetchCategories) createCategory(current *gabs.Container) dto.Category {
	return dto.Category{
		Id:       strconv.Itoa(int(current.Path("id").Data().(float64))),
		Name:     current.Path("name").Data().(string),
		ParentId: strconv.Itoa(int(current.Path("parent_id").Data().(float64))),
		Slug:     current.Path("slug").Data().(string),
		Level:    int(current.Path("level").Data().(float64)),
		Path:     strings.Split(current.Path("path").Data().(string), "/"),
		IsActive: current.Path("is_active").Data().(bool),
	}
}
