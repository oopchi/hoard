package hoard_test

import (
	"fmt"

	"github.com/oopchi/hoard"
)

func ExampleEquipDefault_panic() {
	type nonRegisteredService struct{}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	// Trying to equip a service that wasn't hoarded, this will panic
	hoard.EquipDefault[nonRegisteredService]()
	// Output: Recovered from panic: interface conversion: interface {} is nil, not hoard_test.nonRegisteredService
}
