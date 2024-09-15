package hoard

import (
	"reflect"
	"strings"
	"sync"
)

const (
	// defaultInventoryName is the default Inventory used if no custom name is provided.
	defaultInventoryName = "default"
)

// hoardConfig is a struct that holds the configuration to be used when calling the [Hoard] function.
// This struct is used internally and should not be used directly.
// To specify the desired configuration, use the [HoardOptions] type when calling the [Hoard] function instead.
type hoardConfig struct {

	// shouldReplaceGlobal is a boolean that determines whether the global hoarder should be replaced with the new hoarder created by the [Hoard] function.
	// By default, it is set to true.
	shouldReplaceGlobal bool

	// customHoarder is a custom hoarder that can be used to be merged when calling the [Hoard] function.
	customHoarder Hoarder
}

var (
	// defaultHoardConfig is the default configuration to be used when calling the [Hoard] function.
	// Override the default configuration by specifying the desired configuration in the [HoardOptions] each time when calling the [Hoard] function.
	defaultHoardConfig = hoardConfig{
		shouldReplaceGlobal: true,
	}
)

// HoardOptions is a type that holds the options to be used when calling the [Hoard] function.
// Specifying the desired options in the [HoardOptions] when calling the [Hoard] function will override the default configuration.
// Example usage:
//
//	Hoard(HoardOptions{}.ShouldReplaceGlobal(true), things...)
type HoardOptions []*funcHoardOptions

// funcHoardOptions is a struct that holds a function that modifies the [hoardConfig] struct.
// The function is used to apply the desired configuration to the [hoardConfig] struct.
// This struct is used in the [HoardOptions] type internally and should not be used directly.
// To specify the desired configuration, use the [HoardOptions] type when calling the [Hoard] function instead.
type funcHoardOptions struct {
	f func(*hoardConfig) *hoardConfig
}

// apply is a method that applies a side effect to the [hoardConfig] struct using the function stored in the [funcHoardOptions] struct.
// The method returns the modified [hoardConfig] struct.
// This method is used internally and should not be used directly.
// To specify the desired configuration, use the [HoardOptions] type when calling the [Hoard] function instead.
func (fho *funcHoardOptions) apply(ho *hoardConfig) *hoardConfig {
	return fho.f(ho)
}

// newFuncHoardOptions is a function that creates a new [funcHoardOptions] struct with the given function.
// This function is used internally and should not be used directly.
// To specify the desired configuration, use the [HoardOptions] type when calling the [Hoard] function instead.
func newFuncHoardOptions(f func(*hoardConfig) *hoardConfig) *funcHoardOptions {
	return &funcHoardOptions{f: f}
}

// ShouldReplaceGlobal is a method that sets the [shouldReplaceGlobal] field in the [hoardConfig] struct to the given value.
// The method returns a new [HoardOptions] with the updated configuration.
// Typical usage of this method is to specify whether the global hoarder should be replaced with the new hoarder created by the [Hoard] function.
// If the global hoarder should be replaced, set the value to true. Otherwise, set the value to false.
// The default value is true if this method is not called.
// Example usage:
//
//	Hoard(HoardOptions{}.ShouldReplaceGlobal(true), things...)
func (h HoardOptions) ShouldReplaceGlobal(shouldReplaceGlobal bool) HoardOptions {
	return append(h, newFuncHoardOptions(func(opt *hoardConfig) *hoardConfig {
		opt.shouldReplaceGlobal = shouldReplaceGlobal
		return opt
	}))
}

// WithCustomHoarder is a method that sets the [customHoarder] field in the [hoardConfig] struct to the given value.
// The method returns a new [HoardOptions] with the updated configuration.
// Typical usage of this method is to specify a custom hoarder to be merged when calling the [Hoard] function.
// Example usage:
//
//	Hoard(HoardOptions{}.WithCustomHoarder(customHoarder), things...)
func (h HoardOptions) WithCustomHoarder(customHoarder Hoarder) HoardOptions {
	return append(h, newFuncHoardOptions(func(opt *hoardConfig) *hoardConfig {
		opt.customHoarder = customHoarder
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
	loadout() func(func(string, Inventory) bool)
	merge(hoarder Hoarder)
}

var (
	once          sync.Once = sync.Once{}
	globalHoarder Hoarder
)

type hoarder struct {
	inventoryMap map[string]Inventory

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

	if cfg.customHoarder != nil {
		cfg.customHoarder.merge(h)
		return cfg.customHoarder
	}

	return h
}

func RememberAs(thing interface{}, name string) Item {
	typeOfThing := getTypeOfThing(thing)
	thingName := getCustomThingName(name, typeOfThing)
	return newItem(thing, thingName)
}

func UseInventory(name string) Inventory {
	name = getCustomInventoryName(name)
	inventoryImpl := newInventory(name)

	return inventoryImpl
}

// EquipDefault is a convenience function that returns the requested thing from the default inventoryImpl.
func EquipDefault[T any](customHoarder ...Hoarder) T {
	return EquipWithOption[T](nil, customHoarder...)
}

// EquipFrom is a convenience function that returns the requested thing from the specified inventoryImpl.
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

	inventoryImpl := h.inventoryMap[inventoryName]

	thingName := getCustomThingName(itemName, typeOfThing)

	if v := inventoryImpl.equip(thingName); v != nil {
		return v.use()
	}

	aliasName := getAliasThingName(thingName)

	if aliasName != "" {
		if v := inventoryImpl.equip(aliasName); v != nil {
			return v.use()
		}
	}

	if typeOfThing.Kind() == reflect.Interface {
		for _, thing := range inventoryImpl.loadout() {
			if reflect.TypeOf(thing.use()).Implements(typeOfThing) {
				return thing.use()
			}
		}
	}

	return nil
}

func (h *hoarder) loadout() func(func(string, Inventory) bool) {
	return func(yield func(string, Inventory) bool) {
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
			h.inventoryMap = make(map[string]Inventory)
		}

		if _, ok := h.inventoryMap[k]; !ok {
			h.inventoryMap[k] = v
			continue
		}

		h.inventoryMap[k].merge(v)
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
	inventoryMap := make(map[string]Inventory)
	inventoryMap[defaultInventoryName] = newInventory(defaultInventoryName)
	for _, thing := range things {
		if thing == nil {
			continue
		}

		if v, ok := thing.(Inventory); ok {
			inventoryMap[v.getName()] = newInventory(v.getName())

			// also Put the inventoryImpl items into the default inventoryImpl if absent
			for _, itemImpl := range v.loadout() {
				inventoryMap[defaultInventoryName].
					PutIfAbsent(
						newItem(
							itemImpl.use(),
							getOriginalThingName(itemImpl.getName()),
						),
					)
				inventoryMap[v.getName()].
					Put(
						newItem(
							itemImpl.use(),
							getOriginalThingName(itemImpl.getName()),
						),
					)

				if getAliasThingName(itemImpl.getName()) == "" {
					continue
				}

				inventoryMap[defaultInventoryName].
					PutIfAbsent(
						newItem(
							itemImpl.use(),
							getAliasThingName(itemImpl.getName()),
						),
					).
					PutIfAbsent(
						newItem(
							itemImpl.use(),
							itemImpl.getName(),
						),
					)

				inventoryMap[v.getName()].
					Put(
						newItem(
							itemImpl.use(),
							getAliasThingName(itemImpl.getName()),
						),
					).
					Put(
						newItem(
							itemImpl.use(),
							itemImpl.getName(),
						),
					)
			}
			continue
		}

		if v, ok := thing.(Item); ok {
			inventoryMap[defaultInventoryName].
				PutIfAbsent(
					newItem(
						v.use(),
						getOriginalThingName(v.getName()),
					),
				)

			if getAliasThingName(v.getName()) == "" {
				continue
			}

			inventoryMap[defaultInventoryName].
				Put(
					newItem(
						v.use(),
						getAliasThingName(v.getName()),
					),
				).
				Put(
					newItem(
						v.use(),
						v.getName(),
					),
				)
			continue
		}

		typeOfThing := getTypeOfThing(thing)
		thingName := getThingName(typeOfThing)

		if thingName == "" {
			continue
		}

		inventoryMap[defaultInventoryName].Put(newItem(thing, thingName))
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
	return customInventoryName + defaultInventoryName
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
