package hoard

type Item interface {
	getName() string
	use() interface{}
}

func newItem(thing interface{}, name string) Item {
	return &itemImpl{
		thing: thing,
		name:  name,
	}
}

type itemImpl struct {
	thing interface{}
	name  string
}

func (i *itemImpl) getName() string {
	return i.name
}

func (i *itemImpl) use() interface{} {
	return i.thing
}
