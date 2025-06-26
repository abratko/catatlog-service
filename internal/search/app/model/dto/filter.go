package dto

type Filter struct {
	Key           string `json:"key"`
	Label         string `json:"label"`
	Field         string `json:"field"`
	OperationType string `json:"operationType"`
	Value         any    `json:"value"`
	DefaultValue  any    `json:"defaultValue"`
}
