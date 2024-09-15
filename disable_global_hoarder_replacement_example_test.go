package hoard_test

import (
	"fmt"

	"github.com/oopchi/hoard"
)

func ExampleHoard_withOption() {
	// Hoard with the option to disable global hoarder replacement
	customHoarder := hoard.Hoard(hoard.HoardOptions{}.ShouldReplaceGlobal(false), 42)
	customHoarder = hoard.Hoard(hoard.HoardOptions{}.ShouldReplaceGlobal(false).WithCustomHoarder(customHoarder), 42)

	// Hoard with global hoarder replaced
	hoard.Hoard(nil, 50)
	hoard.Hoard(nil, 70)

	// Equip from the custom hoarder without affecting the global one
	myInt := hoard.EquipDefault[int](customHoarder)

	// Equip from default hoarder will search for items within the global hoarder
	myInt2 := hoard.EquipDefault[int]()

	fmt.Println(myInt, myInt2) // Output: 42 70
}
