package hoard_test

import (
	"fmt"

	"github.com/oopchi/hoard"
)

type EService interface {
	Execute() string
}

type EServiceA struct{}

func (s EServiceA) Execute() string { return "Service A" }

type EServiceB struct{}

func (s EServiceB) Execute() string { return "Service B" }

func ExampleEquipWithOption_interfaces() {
	// Hoard services and annotate them
	hoard.Hoard(nil, hoard.RememberAs(EServiceA{}, "serviceA"))
	hoard.Hoard(nil, hoard.RememberAs(EServiceB{}, "serviceB"))

	// Equip services by their annotation
	sA := hoard.EquipWithOption[EService](hoard.EquipOptions{}.WithCustomItemName("serviceA"))
	sB := hoard.EquipWithOption[EService](hoard.EquipOptions{}.WithCustomItemName("serviceB"))

	// Execute the services
	fmt.Println(sA.Execute(), sB.Execute())
	// Output: Service A Service B
}
