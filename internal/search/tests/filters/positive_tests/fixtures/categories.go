package fixtures

import "gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/category/app/dto"

var CategoriesFixture = []dto.Category{
	{
		Id:   "48",
		Slug: "slug1",
		Name: "Category 1",
		Path: []string{"path1"},
	},
	{
		Id:   "195",
		Slug: "slug2",
		Name: "Category 2",
		Path: []string{"path2"},
	},
	{
		Id:   "354",
		Slug: "slug3",
		Name: "Category 3",
		Path: []string{"path3"},
	},
	{
		Id:   "110",
		Slug: "slug4",
		Name: "Category 4",
		Path: []string{"path4"},
	},
	{
		Id:   "112",
		Slug: "slug5",
		Name: "Category 5",
		Path: []string{"path5"},
	},
}
