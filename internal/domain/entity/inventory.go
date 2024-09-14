package entity

import (
	"slices"
	"sync"
)

type Inventory interface {
	GetName() string

	// Equip returns the item with the given name.
	// Should only be used internally.
	// Prefer using [EquipDefault] or [EquipWithOption] instead.
	Equip(name string) Item
	Put(item Item) Inventory
	PutIfAbsent(item Item) Inventory
	Merge(inventory Inventory) Inventory
	Clone() Inventory

	Loadout() func(func(string, Item) bool)
}

func NewInventory(name string) Inventory {
	return &inventory{
		sortedKeys: make([]string, 0),
		itemMap:    make(map[string]Item),
		name:       name,
		mu:         sync.RWMutex{},
	}
}

type inventory struct {
	sortedKeys []string
	itemMap    map[string]Item
	name       string
	mu         sync.RWMutex
}

func (b *inventory) Clone() Inventory {
	b.mu.RLock()
	defer b.mu.RUnlock()

	clone := NewInventory(b.name)

	for _, k := range b.sortedKeys {
		clone.Put(b.itemMap[k])
	}

	return clone
}

func (b *inventory) GetName() string {
	return b.name
}

func (b *inventory) Equip(name string) Item {
	b.mu.RLock()
	defer b.mu.RUnlock()

	v, ok := b.itemMap[name]

	if !ok {
		return nil
	}

	return v
}

func (b *inventory) Put(item Item) Inventory {
	if item == nil {
		return b
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.itemMap[item.GetName()] = item
	b.sortedKeys = slices.DeleteFunc(b.sortedKeys, func(e string) bool {
		return e == item.GetName()
	})
	b.sortedKeys = append(b.sortedKeys, item.GetName())

	return b
}

func (b *inventory) PutIfAbsent(item Item) Inventory {
	if item == nil {
		return b
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.itemMap[item.GetName()]; !ok {
		b.itemMap[item.GetName()] = item
		b.sortedKeys = append(b.sortedKeys, item.GetName())
	}

	return b
}

func (b *inventory) Merge(inventory Inventory) Inventory {
	if inventory == nil {
		return b
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	for k, v := range inventory.Loadout() {
		b.itemMap[k] = v

		// re-insert the key to ensure the order is consistent
		b.sortedKeys = slices.DeleteFunc(b.sortedKeys, func(e string) bool {
			return e == k
		})

		b.sortedKeys = append(b.sortedKeys, k)
	}

	return b
}

func (b *inventory) Loadout() func(func(string, Item) bool) {
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
