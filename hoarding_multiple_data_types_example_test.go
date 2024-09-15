package hoard_test

import (
	"fmt"

	"github.com/oopchi/hoard"
)

type MyStruct struct {
	Name string
}

func ExampleHoard_multipleDataTypes() {
	// Hoard services of different types
	hoard.Hoard(nil, 123, "test", true, MyStruct{Name: "hoard"}, &MyStruct{Name: "pointer"})

	// Equip services
	i := hoard.EquipDefault[int]()
	s := hoard.EquipDefault[string]()
	b := hoard.EquipDefault[bool]()
	st := hoard.EquipDefault[MyStruct]()
	p := hoard.EquipDefault[*MyStruct]()

	fmt.Println(i, s, b, st.Name, p.Name)
	// Output: 123 test true hoard pointer
}
