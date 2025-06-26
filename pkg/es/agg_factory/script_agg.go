package agg_factory

type ScriptAgg struct {
	Script string `json:"script"`
}

func NewScriptAgg(script string, operation string) m {
	return m{
		operation: m{
			"script": script,
		},
	}
}
