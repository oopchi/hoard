package hoard

import (
	"slices"
	"sync"
)

type inventory interface {
	getName() string

	// equip returns the itemImpl with the given name.
	// Should only be used internally.
	// Prefer using [EquipDefault] or [EquipWithOption] instead.
	equip(name string) item
	put(itemImpl item) inventory
	putIfAbsent(itemImpl item) inventory
	merge(inventoryImpl inventory) inventory

	loadout() func(func(string, item) bool)
}

func newInventory(name string) inventory {
	return &inventoryImpl{
		sortedKeys: make([]string, 0),
		itemMap:    make(map[string]item),
		name:       name,
		mu:         sync.RWMutex{},
	}
}

type inventoryImpl struct {
	sortedKeys []string
	itemMap    map[string]item
	name       string
	mu         sync.RWMutex
}

func (b *inventoryImpl) getName() string {
	return b.name
}

func (b *inventoryImpl) equip(name string) item {
	b.mu.RLock()
	defer b.mu.RUnlock()

	v, ok := b.itemMap[name]

	if !ok {
		return nil
	}

	return v
}

func (b *inventoryImpl) put(itemImpl item) inventory {
	if itemImpl == nil {
		return b
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.itemMap[itemImpl.getName()] = itemImpl
	b.sortedKeys = slices.DeleteFunc(b.sortedKeys, func(e string) bool {
		return e == itemImpl.getName()
	})
	b.sortedKeys = append(b.sortedKeys, itemImpl.getName())

	return b
}

func (b *inventoryImpl) putIfAbsent(itemImpl item) inventory {
	if itemImpl == nil {
		return b
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.itemMap[itemImpl.getName()]; !ok {
		b.itemMap[itemImpl.getName()] = itemImpl
		b.sortedKeys = append(b.sortedKeys, itemImpl.getName())
	}

	return b
}

func (b *inventoryImpl) merge(inventoryImpl inventory) inventory {
	if inventoryImpl == nil {
		return b
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	for k, v := range inventoryImpl.loadout() {
		b.itemMap[k] = v

		// re-insert the key to ensure the order is consistent
		b.sortedKeys = slices.DeleteFunc(b.sortedKeys, func(e string) bool {
			return e == k
		})

		b.sortedKeys = append(b.sortedKeys, k)
	}

	return b
}

func (b *inventoryImpl) loadout() func(func(string, item) bool) {
	return func(yield func(string, item) bool) {
		b.mu.RLock()
		defer b.mu.RUnlock()

		for _, k := range b.sortedKeys {
			if !yield(k, b.itemMap[k]) {
				break
			}
		}
	}
}
