package hoard

import (
	"reflect"
	"strings"
	"sync"

	"github.com/oopchi/hoard/internal/core/constant"
	"github.com/oopchi/hoard/internal/domain/entity"
)

type hoardConfig struct {
	shouldReplaceGlobal bool
}

var (
	defaultHoardConfig = hoardConfig{
		shouldReplaceGlobal: true,
	}
)

type HoardOptions []*funcHoardOptions

type funcHoardOptions struct {
	f func(*hoardConfig) *hoardConfig
}

func (fho *funcHoardOptions) apply(ho *hoardConfig) *hoardConfig {
	return fho.f(ho)
}

func newFuncHoardOptions(f func(*hoardConfig) *hoardConfig) *funcHoardOptions {
	return &funcHoardOptions{f: f}
}

func (h HoardOptions) ShouldReplaceGlobal(shouldReplaceGlobal bool) HoardOptions {
	return append(h, newFuncHoardOptions(func(opt *hoardConfig) *hoardConfig {
		opt.shouldReplaceGlobal = shouldReplaceGlobal
		return opt
	}))
}

type equipConfig struct {
	customInventoryName string
	customItemName      string
}

var (
	defaultEquipConfig = equipConfig{}
)

type EquipOptions []*funcEquipOptions

type funcEquipOptions struct {
	f func(*equipConfig) *equipConfig
}

func (fho *funcEquipOptions) apply(ho *equipConfig) *equipConfig {
	return fho.f(ho)
}

func newFuncEquipOptions(f func(*equipConfig) *equipConfig) *funcEquipOptions {
	return &funcEquipOptions{f: f}
}

func (h EquipOptions) WithCustomInventoryName(customInventoryName string) EquipOptions {
	return append(h, newFuncEquipOptions(func(opt *equipConfig) *equipConfig {
		opt.customInventoryName = customInventoryName
		return opt
	}))
}

func (h EquipOptions) WithCustomItemName(customItemName string) EquipOptions {
	return append(h, newFuncEquipOptions(func(opt *equipConfig) *equipConfig {
		opt.customItemName = customItemName
		return opt
	}))
}

type Hoarder interface {
	get(typeOfThing reflect.Type, inventoryName, itemName string) interface{}
	loadout() func(func(string, entity.Inventory) bool)
	merge(hoarder Hoarder)
}

var (
	once          sync.Once = sync.Once{}
	globalHoarder Hoarder
)

type hoarder struct {
	inventoryMap map[string]entity.Inventory

	mu sync.RWMutex
}

func Hoard(opt HoardOptions, things ...interface{}) Hoarder {
	cfg := defaultHoardConfig

	for _, f := range opt {
		f.apply(&cfg)
	}

	h := factory(things...)

	if cfg.shouldReplaceGlobal {
		initGlobalHoarder(h)

		globalHoarder.merge(h)
	}

	return h
}

func RememberAs(thing interface{}, name string) entity.Item {
	typeOfThing := getTypeOfThing(thing)
	thingName := getCustomThingName(name, typeOfThing)
	return entity.NewItem(thing, thingName)
}

func UseInventory(name string) entity.Inventory {
	name = getCustomInventoryName(name)
	inventory := entity.NewInventory(name)

	return inventory
}

// EquipDefault is a convenience function that returns the requested thing from the default inventory.
func EquipDefault[T any](customHoarder ...Hoarder) T {
	return EquipWithOption[T](nil, customHoarder...)
}

// EquipFrom is a convenience function that returns the requested thing from the specified inventory.
func EquipWithOption[T any](opt EquipOptions, customHoarder ...Hoarder) T {
	cfg := defaultEquipConfig

	for _, f := range opt {
		f.apply(&cfg)
	}

	customInventoryName := cfg.customInventoryName
	customItemName := cfg.customItemName

	typeOfType := reflect.TypeFor[T]()

	var hoarder Hoarder

	if len(customHoarder) > 0 && customHoarder[0] != nil {
		hoarder = customHoarder[0]
	} else {
		hoarder = globalFactory()
	}

	inventoryName := getCustomInventoryName(customInventoryName)

	return hoarder.get(typeOfType, inventoryName, customItemName).(T)
}

func (h *hoarder) get(typeOfThing reflect.Type, inventoryName, itemName string) interface{} {
	if typeOfThing.Kind() == reflect.Func {
		// TODO (oopchi): handle function type
		return nil
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	if _, ok := h.inventoryMap[inventoryName]; !ok {
		return nil
	}

	inventory := h.inventoryMap[inventoryName]

	thingName := getCustomThingName(itemName, typeOfThing)

	if v := inventory.Equip(thingName); v != nil {
		return v.Use()
	}

	aliasName := getAliasThingName(thingName)

	if aliasName != "" {
		if v := inventory.Equip(aliasName); v != nil {
			return v.Use()
		}
	}

	if typeOfThing.Kind() == reflect.Interface {
		for _, thing := range inventory.Loadout() {
			if reflect.TypeOf(thing.Use()).Implements(typeOfThing) {
				return thing.Use()
			}
		}
	}

	return nil
}

func (h *hoarder) loadout() func(func(string, entity.Inventory) bool) {
	return func(yield func(string, entity.Inventory) bool) {
		h.mu.RLock()
		defer h.mu.RUnlock()

		for k, v := range h.inventoryMap {
			if !yield(k, v) {
				break
			}
		}
	}
}

func (h *hoarder) merge(hoarder Hoarder) {
	if hoarder == nil {
		return
	}

	if h == hoarder {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for k, v := range hoarder.loadout() {
		if h.inventoryMap == nil {
			h.inventoryMap = make(map[string]entity.Inventory)
		}

		if _, ok := h.inventoryMap[k]; !ok {
			h.inventoryMap[k] = v
			continue
		}

		h.inventoryMap[k].Merge(v)
	}
}

func globalFactory() Hoarder {
	if globalHoarder != nil {
		return globalHoarder
	}

	initGlobalHoarder(factory())

	return globalHoarder
}

func initGlobalHoarder(h Hoarder) {
	once.Do(func() {
		globalHoarder = h
	})
}

func factory(things ...interface{}) Hoarder {
	inventoryMap := make(map[string]entity.Inventory)
	inventoryMap[constant.DEFAULT_INVENTORY_NAME] = entity.NewInventory(constant.DEFAULT_INVENTORY_NAME)
	for _, thing := range things {
		if thing == nil {
			continue
		}

		if v, ok := thing.(entity.Inventory); ok {
			inventoryMap[v.GetName()] = entity.NewInventory(v.GetName())

			// also put the inventory items into the default inventory if absent
			for _, item := range v.Loadout() {
				inventoryMap[constant.DEFAULT_INVENTORY_NAME].
					PutIfAbsent(
						entity.NewItem(
							item.Use(),
							getOriginalThingName(item.GetName()),
						),
					)
				inventoryMap[v.GetName()].
					Put(
						entity.NewItem(
							item.Use(),
							getOriginalThingName(item.GetName()),
						),
					)

				if getAliasThingName(item.GetName()) == "" {
					continue
				}

				inventoryMap[constant.DEFAULT_INVENTORY_NAME].
					PutIfAbsent(
						entity.NewItem(
							item.Use(),
							getAliasThingName(item.GetName()),
						),
					).
					PutIfAbsent(
						entity.NewItem(
							item.Use(),
							item.GetName(),
						),
					)

				inventoryMap[v.GetName()].
					Put(
						entity.NewItem(
							item.Use(),
							getAliasThingName(item.GetName()),
						),
					).
					Put(
						entity.NewItem(
							item.Use(),
							item.GetName(),
						),
					)
			}
			continue
		}

		if v, ok := thing.(entity.Item); ok {
			inventoryMap[constant.DEFAULT_INVENTORY_NAME].
				PutIfAbsent(
					entity.NewItem(
						v.Use(),
						getOriginalThingName(v.GetName()),
					),
				)

			if getAliasThingName(v.GetName()) == "" {
				continue
			}

			inventoryMap[constant.DEFAULT_INVENTORY_NAME].
				Put(
					entity.NewItem(
						v.Use(),
						getAliasThingName(v.GetName()),
					),
				).
				Put(
					entity.NewItem(
						v.Use(),
						v.GetName(),
					),
				)
			continue
		}

		typeOfThing := getTypeOfThing(thing)
		thingName := getThingName(typeOfThing)

		if thingName == "" {
			continue
		}

		inventoryMap[constant.DEFAULT_INVENTORY_NAME].Put(entity.NewItem(thing, thingName))
	}

	return &hoarder{
		mu:           sync.RWMutex{},
		inventoryMap: inventoryMap,
	}
}

func getThingName(typeOfThing reflect.Type) string {
	if typeOfThing == nil {
		return ""
	}

	prefix := ""

	if typeOfThing.Kind() == reflect.Pointer {
		typeOfThing = typeOfThing.Elem()
		prefix = "*"
	}

	if typeOfThing.Kind() == reflect.Func {
		// TODO (oopchi): handle function type
		return ""
	}

	return prefix + typeOfThing.PkgPath() + typeOfThing.Name()
}

func getTypeOfThing(thing interface{}) reflect.Type {
	if thing == nil {
		return nil
	}

	return reflect.TypeOf(thing)
}

func getCustomInventoryName(customInventoryName string) string {
	return customInventoryName + constant.DEFAULT_INVENTORY_NAME
}

func getCustomThingName(customItemName string, typeOfThing reflect.Type) string {
	if customItemName == "" {
		return getThingName(typeOfThing)
	}

	return getThingName(typeOfThing) + "\n" + customItemName
}

func getAliasThingName(thingName string) string {
	if len(strings.Split(thingName, "\n")) <= 1 {
		return ""
	}

	return strings.Join(strings.Split(thingName, "\n")[1:], "\n")
}

func getOriginalThingName(thingName string) string {
	return strings.Split(thingName, "\n")[0]
}
