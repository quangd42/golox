package lox

type array struct {
	value []any
}

func newArray() *array {
	return &array{value: make([]any, 0)}
}

func (a *array) Assign(idx int, val any) {
	a.value[idx] = val
}

func (a *array) Get(idx int) any {
	return a.value[idx]
}

func (a *array) Append(vals ...any) {
	a.value = append(a.value, vals...)
}

func (a *array) Len() int {
	return len(a.value)
}
