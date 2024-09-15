package hoard_test

import (
	"fmt"

	"github.com/oopchi/hoard"
)

type SMyService interface {
	Execute() string
}

type ServiceImpl struct {
	ID int
}

func (s ServiceImpl) Execute() string {
	return fmt.Sprintf("Executing Service with ID: %d\n", s.ID)
}

func ExampleEquipWithOption_inventories() {
	// Use custom Inventory and annotation to hoard multiple implementations
	hoard.Hoard(nil,
		hoard.UseInventory("legendary items inventory").
			Put(hoard.RememberAs(ServiceImpl{ID: 1}, "impl1")).
			Put(hoard.RememberAs(ServiceImpl{ID: 1}, "")). // pass an empty string to use the default item name
			Put(hoard.RememberAs(ServiceImpl{ID: 2}, "impl2")),
	)

	// Equip the services using annotations
	svc1 := hoard.EquipWithOption[SMyService](hoard.EquipOptions{}.WithCustomItemName("impl1").WithCustomInventoryName("legendary items inventory"))

	// You can also skip the inventory name because if an item is just hoarded for the first time (no other same item has been hoarded yet)
	// then no matter the custom inventory used, it will also be stored at the default inventory
	svc2 := hoard.EquipWithOption[SMyService](hoard.EquipOptions{}.WithCustomItemName("impl2"))

	// You can even skip the annotation whatsoever if there has only ever been one such item being hoarded even if its annotated
	// If there were already multiple such items being hoarded though, if its being hoarded through a custom inventory
	// then it won't override the one at the default inventory anymore, however it will still override the one at that custom inventory if any existed
	svc3 := hoard.EquipDefault[ServiceImpl]()

	// You can however override the default inventory again if you specifically hoard on the default inventory (hoarding without specifying [hoard.UseInventory])
	hoard.Hoard(nil, ServiceImpl{ID: 5}, ServiceImpl{ID: 8})

	svc4 := hoard.EquipDefault[ServiceImpl]()

	fmt.Println(svc1.Execute(), svc2.Execute(), svc3.Execute(), svc4.Execute())
	// Output: Executing Service with ID: 1
	//  Executing Service with ID: 2
	//  Executing Service with ID: 1
	//  Executing Service with ID: 8
}
