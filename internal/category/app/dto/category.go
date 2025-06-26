package dto

type Category struct {
	Id       string      `json:"id"`
	Slug     string      `json:"slug"`
	Name     string      `json:"name"`
	Path     []string    `json:"path"`
	IsActive bool        `json:"isActive"`
	Level    int         `json:"level"`
	ParentId string      `json:"parentId"`
	Children []*Category `json:"children"`
}
