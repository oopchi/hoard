package hoard_test

import (
	"fmt"

	"github.com/oopchi/hoard"
)

type Service struct {
	Name string
}

func ExampleHoard_multipleItemsOfTheSameType() {
	// Hoard multiple services of the same type
	hoard.Hoard(nil, hoard.RememberAs(Service{Name: "Service 1"}, "service1"))
	hoard.Hoard(nil, hoard.RememberAs(Service{Name: "Service 2"}, "service2"))

	// Equip services by name
	s1 := hoard.EquipWithOption[Service](hoard.EquipOptions{}.WithCustomItemName("service1"))
	s2 := hoard.EquipWithOption[Service](hoard.EquipOptions{}.WithCustomItemName("service2"))

	fmt.Println(s1.Name, s2.Name)
	// Output: Service 1 Service 2
}
