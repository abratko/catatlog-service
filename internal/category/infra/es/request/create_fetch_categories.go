package request

type M map[string]any

func CreateFetchCategories() M {
	return M{
		"query": M{
			"bool": M{
				"filter": M{
					"terms": M{
						"level": []int{1},
					},
				},
			},
		},
	}
}
