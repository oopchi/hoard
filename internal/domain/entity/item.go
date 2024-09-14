package entity

type Item interface {
	GetName() string
	Use() interface{}
}

func NewItem(thing interface{}, name string) Item {
	return &item{
		thing: thing,
		name:  name,
	}
}

type item struct {
	thing interface{}
	name  string
}

func (i *item) GetName() string {
	return i.name
}

func (i *item) Use() interface{} {
	return i.thing
}
