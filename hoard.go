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

	// customHoarder is a custom hoarder that can be used to be merged and returned when calling the [Hoard] function.
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

// equipConfig is a struct that holds the configuration to be used when calling the [EquipWithOption] function.
// This struct is used internally and should not be used directly.
// To specify the desired configuration, use the [EquipOptions] type when calling the [EquipWithOption] function instead.
type equipConfig struct {

	// customInventoryName is a custom inventory name that can be used to get the desired thing from the specified inventory.
	customInventoryName string

	// customItemName is a custom [Item] name that can be used to get the desired thing from the specified inventory.
	customItemName string
}

var (
	// defaultEquipConfig is the default configuration to be used when calling the [EquipWithOption] function.
	defaultEquipConfig = equipConfig{}
)

// EquipOptions is a type that holds the options to be used when calling the [EquipWithOption] function.
// Specifying the desired options in the [EquipOptions] when calling the [EquipWithOption] function will override the default configuration.
// Example usage:
//
//	EquipWithOption(EquipOptions{}.WithCustomInventoryName("customInventoryName"), customHoarder...)
type EquipOptions []*funcEquipOptions

// funcEquipOptions is a struct that holds a function that modifies the [equipConfig] struct.
// The function is used to apply the desired configuration to the [equipConfig] struct.
// This struct is used in the [EquipOptions] type internally and should not be used directly.
// To specify the desired configuration, use the [EquipOptions] type when calling the [EquipWithOption] function instead.
type funcEquipOptions struct {
	f func(*equipConfig) *equipConfig
}

// apply is a method that applies a side effect to the [equipConfig] struct using the function stored in the [funcEquipOptions] struct.
func (fho *funcEquipOptions) apply(ho *equipConfig) *equipConfig {
	return fho.f(ho)
}

// newFuncEquipOptions is a function that creates a new [funcEquipOptions] struct with the given function.
func newFuncEquipOptions(f func(*equipConfig) *equipConfig) *funcEquipOptions {
	return &funcEquipOptions{f: f}
}

// WithCustomInventoryName is a method that sets the [customInventoryName] field in the [equipConfig] struct to the given value.
// The method returns a new [EquipOptions] with the updated configuration.
// Typical usage of this method is to specify a custom inventory name to get the desired thing from the specified inventory.
// Example usage:
//
//	EquipWithOption(EquipOptions{}.WithCustomInventoryName("customInventoryName"), customHoarder...)
func (h EquipOptions) WithCustomInventoryName(customInventoryName string) EquipOptions {
	return append(h, newFuncEquipOptions(func(opt *equipConfig) *equipConfig {
		opt.customInventoryName = customInventoryName
		return opt
	}))
}

// WithCustomItemName is a method that sets the [customItemName] field in the [equipConfig] struct to the given value.
// The method returns a new [EquipOptions] with the updated configuration.
// Typical usage of this method is to specify a custom item name to get the desired thing from the specified inventory.
// Example usage:
//
//	EquipWithOption(EquipOptions{}.WithCustomItemName("customItemName"), customHoarder...)
func (h EquipOptions) WithCustomItemName(customItemName string) EquipOptions {
	return append(h, newFuncEquipOptions(func(opt *equipConfig) *equipConfig {
		opt.customItemName = customItemName
		return opt
	}))
}

// Hoarder is an interface that defines the methods to be used internally by the [Hoard] function.
// This interface is used internally and should not be used directly.
// To create a new hoarder, use the [Hoard] function instead.
type Hoarder interface {

	// get is a method that returns the requested thing from the specified inventory.
	// The method returns the requested thing if found. Otherwise, it returns nil.
	// This method is used internally and should not be used directly.
	// This method is thread-safe.
	get(typeOfThing reflect.Type, inventoryName, itemName string) interface{}

	// loadout is a method that returns the inventory map.
	// The method returns the inventory map.
	// This method is used internally and should not be used directly.
	// This method is thread-safe.
	loadout() func(func(string, Inventory) bool)

	// merge is a method that merges the given hoarder with the current hoarder.
	// This method is used internally and should not be used directly.
	// This method is thread-safe.
	merge(hoarder Hoarder)
}

var (
	once sync.Once = sync.Once{}

	// globalHoarder is the global hoarder that holds the inventory map.
	globalHoarder Hoarder
)

// hoarder is a struct that implements the [Hoarder] interface.
// This struct is used internally and should not be used directly.
// To create a new hoarder, use the [Hoard] function instead.
// All methods in this struct are thread-safe.
type hoarder struct {

	// inventoryMap is a map that holds multiple inventories for managing items.
	inventoryMap map[string]Inventory

	mu sync.RWMutex
}

// Hoard is a function that creates a new [Hoarder] with the given things and options.
// The function returns a new [Hoarder] which can be used as a custom [Hoarder] or ignored.
//
// By default, registering things with the [Hoard] function will replace existing things of the same type and alias inside the global hoarder.
// To disable global [Hoarder] replacement, use the [HoardOptions] to specify the desired configuration.
// By default, the global [Hoarder] is the only one being merged with the new [Hoarder] created by the [Hoard] function.
// You can specify a custom [Hoarder] to be merged with the new [Hoarder] by using the [HoardOptions] to specify the desired configuration.
//
// Multiple things can be registered at once by providing a list of things.
// Registering multiple things with the same type and alias will replace existing things of the same type and alias inside the global [Hoarder].
//
// To register the thing with a custom name, wrap the thing with the [RememberAs] function.
//
// To register the thing with a custom inventory, insert the thing with the [Inventory.Put] method obtained from the [UseInventory] function and pass the returned [Inventory] to this function.
//
// The [Hoard] function is thread-safe.
//
// Example usage:
//
//	Hoard(nil, 42)
//	Hoard(HoardOptions{}.ShouldReplaceGlobal(false), 42)
//	Hoard(HoardOptions{}.ShouldReplaceGlobal(false).WithCustomHoarder(customHoarder), 42)
//	Hoard(nil, RememberAs(42, "customName"))
//	Hoard(nil, UseInventory("customInventory").Put(RememberAs(42, "customName")))
//	Hoard(nil, UseInventory("customInventory").Put(RememberAs(42, "")))
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

// RememberAs is a function that wraps the given thing with a custom name.
// The function returns a new [Item] with the given thing and custom name.
// The custom name is used to identify the thing when registering it with the [Hoard] function.
// Passing an empty string as the custom name will use the default name.
// Example usage:
//
//	RememberAs(42, "customName")
//	RememberAs(42, "")
func RememberAs(thing interface{}, name string) Item {
	typeOfThing := getTypeOfThing(thing)
	thingName := getCustomThingName(name, typeOfThing)
	return newItem(thing, thingName)
}

// UseInventory is a function that creates a new inventory with the given name.
// The function returns a new [Inventory] with the given name.
// The name is used to identify the inventory when registering items with the [Hoard] function.
// Example usage:
//
//	UseInventory("customInventory").Put(RememberAs(42, "customName")).Put(RememberAs(42, ""))
func UseInventory(name string) Inventory {
	name = getCustomInventoryName(name)
	inventoryImpl := newInventory(name)

	return inventoryImpl
}

// EquipDefault is a convenience function that is equivalent to calling the [EquipWithOption] function with the default configuration or nil [EquipOptions].
func EquipDefault[T any](customHoarder ...Hoarder) T {
	return EquipWithOption[T](nil, customHoarder...)
}

// EquipWithOption is a function that returns the requested thing from the specified [Inventory].
// By default, the thing is retrieved from the default [Inventory].
// The function refers to the global [Hoarder] to get the desired thing unless a custom [Hoarder] is specified.
// The function returns the requested thing if found. Otherwise, it panics.
//
// To specify custom [Item] name or custom [Inventory] name, use the [EquipOptions] when calling the [EquipWithOption] function.
//
// To specify a custom [Hoarder] to be used, pass the custom [Hoarder] as an argument when calling the [EquipWithOption] function.
//
// The [EquipWithOption] function is thread-safe.
//
// Example usage:
//
//	EquipWithOption[Egg](EquipOptions{}.WithCustomInventoryName("customInventoryName").WithCustomItemName("customItemName"))
//	EquipWithOption[Egg](EquipOptions{}.WithCustomInventoryName("customInventoryName"), customHoarder)
//	EquipWithOption[Egg](nil, customHoarder)
//	EquipWithOption[Egg](nil)
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
