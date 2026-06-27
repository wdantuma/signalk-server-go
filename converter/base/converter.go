package base

type InputFields map[string]interface{}

func (field InputFields) Contains(key string) bool {
	_, ok := field[key]
	return ok
}

type Field struct {
	Node    string
	Source  string
	Value   func(InputFields) interface{}
	Filter  func(InputFields) bool
	Context func(InputFields) *string
}
