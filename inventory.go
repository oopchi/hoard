package hoard

import (
	"slices"
	"sync"
)

// Inventory is a collection of items.
// It is used to store items that can be equipped.
// This interface is used internally by the [Hoard], [EquipDefault], and [EquipWithOption] functions.
// To create a custom inventory, use the [UseInventory] function.
type Inventory interface {

	// Put adds an [Item] to the inventory.
	// To get an [Item], refer to the [UseInventory] function.
	//
	// Example:
	// 	Hoard(nil, UseInventory("test").Put(RememberAs("test", "test")))
	Put(item Item) Inventory

	// PutIfAbsent adds an [Item] to the inventory if it does not exist yet.
	// The comparison is done by also comparing the alias of the item, hence if the alias is different, it will be considered a different item.
	// To get an [Item], refer to the [UseInventory] function.
	//
	// Example:
	// 	Hoard(nil, UseInventory("test").PutIfAbsent(RememberAs("test", "test")))
	PutIfAbsent(item Item) Inventory

	getName() string

	// equip returns the itemImpl with the given name.
	// Should only be used internally.
	// Prefer using [EquipDefault] or [EquipWithOption] instead.
	equip(name string) Item

	merge(inventoryImpl Inventory) Inventory

	loadout() func(func(string, Item) bool)
}

func newInventory(name string) Inventory {
	return &inventoryImpl{
		sortedKeys: make([]string, 0),
		itemMap:    make(map[string]Item),
		name:       name,
		mu:         sync.RWMutex{},
	}
}

type inventoryImpl struct {
	sortedKeys []string
	itemMap    map[string]Item
	name       string
	mu         sync.RWMutex
}

// Put adds an [Item] to the inventory.
// To get an [Item], refer to the [UseInventory] function.
//
// Example:
//
//	Hoard(nil, UseInventory("test").Put(RememberAs("test", "test")))
func (b *inventoryImpl) Put(item Item) Inventory {
	if item == nil {
		return b
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.itemMap[item.getName()] = item
	b.sortedKeys = slices.DeleteFunc(b.sortedKeys, func(e string) bool {
		return e == item.getName()
	})
	b.sortedKeys = append(b.sortedKeys, item.getName())

	return b
}

// PutIfAbsent adds an [Item] to the inventory if it does not exist yet.
// The comparison is done by also comparing the alias of the item, hence if the alias is different, it will be considered adifferent item.
// To get an [Item], refer to the [UseInventory] function.
//
// Example:
//
//	Hoard(nil, UseInventory("test").PutIfAbsent(RememberAs("test", "test")))
func (b *inventoryImpl) PutIfAbsent(item Item) Inventory {
	if item == nil {
		return b
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.itemMap[item.getName()]; !ok {
		b.itemMap[item.getName()] = item
		b.sortedKeys = append(b.sortedKeys, item.getName())
	}

	return b
}

func (b *inventoryImpl) getName() string {
	return b.name
}

func (b *inventoryImpl) equip(name string) Item {
	b.mu.RLock()
	defer b.mu.RUnlock()

	v, ok := b.itemMap[name]

	if !ok {
		return nil
	}

	return v
}

func (b *inventoryImpl) merge(invent Inventory) Inventory {
	if invent == nil {
		return b
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	for k, v := range invent.loadout() {
		b.itemMap[k] = v

		// re-insert the key to ensure the order is consistent
		b.sortedKeys = slices.DeleteFunc(b.sortedKeys, func(e string) bool {
			return e == k
		})

		b.sortedKeys = append(b.sortedKeys, k)
	}

	return b
}

func (b *inventoryImpl) loadout() func(func(string, Item) bool) {
	return func(yield func(string, Item) bool) {
		b.mu.RLock()
		defer b.mu.RUnlock()

		for _, k := range b.sortedKeys {
			if !yield(k, b.itemMap[k]) {
				break
			}
		}
	}
}
