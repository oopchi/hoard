package hoard_test

import (
	"fmt"

	"github.com/oopchi/hoard"
)

func ExampleEquipWithOption() {
	// Disable global hoarder replacement
	options := hoard.HoardOptions{}.ShouldReplaceGlobal(false)

	// Hoard services with a custom hoarder
	customHoarder := hoard.Hoard(options, 100, "custom string")

	// Equip services from the custom hoarder
	i := hoard.EquipWithOption[int](nil, customHoarder)
	s := hoard.EquipWithOption[string](nil, customHoarder)

	fmt.Println(i, s) // Output: 100 custom string
}
