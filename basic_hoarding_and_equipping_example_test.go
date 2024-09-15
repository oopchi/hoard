package hoard_test

import (
	"fmt"

	"github.com/oopchi/hoard"
)

type MyService struct {
	Name string
}

func ExampleEquipDefault() {
	// Hoard items of different types
	hoard.Hoard(nil, 42, "Hoarding a string", true, &MyService{Name: "Service 1"})

	// Equip the items
	myInt := hoard.EquipDefault[int]()
	myString := hoard.EquipDefault[string]()
	myBool := hoard.EquipDefault[bool]()
	myService := hoard.EquipDefault[*MyService]()

	fmt.Println(myInt, myString, myBool, myService.Name)
	// Output: 42 Hoarding a string true Service 1
}
