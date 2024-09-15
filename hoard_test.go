package hoard

import (
	"reflect"
	"sync"

	"github.com/stretchr/testify/require"
)

type TestFooImpl struct {
	Name string
}

func (t TestFooImpl) getName() string {
	return t.Name
}

type TestFooer interface {
	getName() string
}

func (s *suiteTest) Test_getTypeOfThing() {

	tests := []struct {
		name  string
		given interface{}
		want  reflect.Type
	}{
		{
			name:  "should be able to get the type of a struct",
			given: TestFooImpl{},
			want:  reflect.TypeOf(TestFooImpl{}),
		},
		{
			name: "should be able to get the type of an interface",
			given: func() TestFooer {
				return TestFooImpl{}
			}(),
			want: reflect.TypeOf(TestFooImpl{}),
		},
		{
			name:  "should be able to get the type of a pointer",
			given: &TestFooImpl{},
			want:  reflect.TypeOf(&TestFooImpl{}),
		},
		{
			name:  "should be able to handle nil by returning nil",
			given: nil,
			want:  nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got := getTypeOfThing(tt.given)

			if tt.want == nil {
				require.Nil(s.T(), got)
				return
			}

			require.Equal(s.T(), tt.want.Kind(), got.Kind())
		})
	}
}

func (s *suiteTest) Test_getThingName() {
	type TestFoo struct {
		Name string
	}

	type testFoo struct {
		privateName string
	}

	tests := []struct {
		name  string
		given reflect.Type
		want  string
	}{
		{
			name:  "should be able to get the name of a struct",
			given: reflect.TypeOf(TestFoo{}),
			want:  "github.com/oopchi/hoardTestFoo",
		},
		{
			name:  "should be able to get the name of an interface",
			given: reflect.TypeOf((*TestFooer)(nil)).Elem(),
			want:  "github.com/oopchi/hoardTestFooer",
		},
		{
			name:  "should be able to handle nil by returning empty string",
			given: nil,
			want:  "",
		},
		{
			name:  "should be able to get the name of a pointer",
			given: reflect.TypeOf(&TestFoo{}),
			want:  "*github.com/oopchi/hoardTestFoo",
		},
		{
			name: "should be able to get the name of a private struct with private fields",
			given: reflect.TypeOf(testFoo{
				privateName: "test",
			}),
			want: "github.com/oopchi/hoardtestFoo",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got := getThingName(tt.given)

			require.Equal(s.T(), tt.want, got)
		})
	}
}

func (s *suiteTest) Test_factory() {
	tests := []struct {
		name  string
		given []interface{}
		want  Hoarder
	}{
		{
			name:  "should be able to create a hoarder without any items",
			given: nil,
			want: &hoarder{
				inventoryMap: map[string]Inventory{
					defaultInventoryName: newInventory(defaultInventoryName),
				},
				mu: sync.RWMutex{},
			},
		},
		{
			name: "should be able to create a hoarder with a single named item",
			given: []interface{}{
				newItem("test thing", "test"),
			},
			want: &hoarder{
				inventoryMap: func() map[string]Inventory {
					inventoryMap := make(map[string]Inventory)
					inventoryMap[defaultInventoryName] = newInventory(defaultInventoryName)
					inventoryMap[defaultInventoryName].Put(newItem("test thing", "test"))

					return inventoryMap
				}(),
				mu: sync.RWMutex{},
			},
		},
		{
			name: "should be able to create a hoarder with a single named inventory",
			given: []interface{}{
				newInventory("test"),
			},
			want: &hoarder{
				inventoryMap: func() map[string]Inventory {
					inventoryMap := make(map[string]Inventory)
					inventoryMap["test"] = newInventory("test")
					inventoryMap[defaultInventoryName] = newInventory(defaultInventoryName)

					return inventoryMap
				}(),
				mu: sync.RWMutex{},
			},
		},
		{
			name: "should be able to create a hoarder with a single named inventory with a single pointer item and a single named item",
			given: []interface{}{
				func() Inventory {
					i := newInventory("test")
					i.Put(newItem(&TestFooImpl{}, "test234"))

					return i
				}(),
				newItem("test thing", "test"),
			},
			want: &hoarder{
				inventoryMap: func() map[string]Inventory {
					inventoryMap := make(map[string]Inventory)
					inventoryMap["test"] = newInventory("test")
					inventoryMap["test"].Put(newItem(&TestFooImpl{}, "test234"))
					inventoryMap[defaultInventoryName] = newInventory(defaultInventoryName).
						Put(newItem("test thing", "test")).
						Put(newItem(&TestFooImpl{}, "test234"))

					return inventoryMap
				}(),
				mu: sync.RWMutex{},
			},
		},
		{
			name: "should be able to create a hoarder with a single named inventory with a single pointer item and a single named item and a single interface item and a nil item and a struct item and a pointer item and a hoarder itself",
			given: []interface{}{
				func() Inventory {
					i := newInventory("test")
					i.Put(newItem(&TestFooImpl{}, "test234"))

					return i
				}(),
				newItem("test thing", "test"),
				func() TestFooer {
					return TestFooImpl{}
				},
				"hehe",
				nil,
				TestFooImpl{},
				&TestFooImpl{},
				&hoarder{
					inventoryMap: map[string]Inventory{
						"test": newInventory("test"),
						"test2": func() Inventory {
							i := newInventory("test2")
							i.Put(newItem("test thing234", "test234"))

							return i
						}(),
					},
					mu: sync.RWMutex{},
				},
			},
			want: &hoarder{
				inventoryMap: func() map[string]Inventory {
					inventoryMap := make(map[string]Inventory)
					inventoryMap["test"] = newInventory("test")
					inventoryMap["test"].Put(newItem(&TestFooImpl{}, "test234"))
					inventoryMap[defaultInventoryName] = newInventory(defaultInventoryName)
					inventoryMap[defaultInventoryName].Put(newItem("test thing", "test"))
					inventoryMap[defaultInventoryName].Put(newItem(TestFooImpl{}, "github.com/oopchi/hoardTestFooImpl"))
					inventoryMap[defaultInventoryName].Put(newItem(&TestFooImpl{}, "*github.com/oopchi/hoardTestFooImpl"))
					inventoryMap[defaultInventoryName].Put(newItem("hehe", "string"))
					inventoryMap[defaultInventoryName].Put(newItem(&hoarder{
						inventoryMap: map[string]Inventory{
							"test": newInventory("test"),
							"test2": func() Inventory {
								i := newInventory("test2")
								i.Put(newItem("test thing234", "test234"))

								return i
							}(),
						},
						mu: sync.RWMutex{},
					}, "*github.com/oopchi/hoardhoarder")).Put(newItem(&TestFooImpl{}, "test234"))

					return inventoryMap
				}(),
				mu: sync.RWMutex{},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got := factory(tt.given...)

			require.NotNil(s.T(), got)

			gotInventories := []Inventory{}
			gotInventoriesNames := []string{}
			for _, v := range got.loadout() {
				gotInventories = append(gotInventories, v)
				gotInventoriesNames = append(gotInventoriesNames, v.getName())
			}

			wantInventories := []Inventory{}
			wantInventoriesNames := []string{}
			for _, v := range tt.want.loadout() {
				wantInventories = append(wantInventories, v)
				wantInventoriesNames = append(wantInventoriesNames, v.getName())
			}

			require.ElementsMatch(s.T(), wantInventoriesNames, gotInventoriesNames)

			gotItems := map[string][]Item{}
			for _, v := range gotInventories {
				for _, it := range v.loadout() {
					if gotItems[v.getName()] == nil {
						gotItems[v.getName()] = []Item{}
					}
					gotItems[v.getName()] = append(gotItems[v.getName()], it)
				}
			}

			wantItems := map[string][]Item{}
			for _, v := range wantInventories {
				for _, it := range v.loadout() {
					if wantItems[v.getName()] == nil {
						wantItems[v.getName()] = []Item{}
					}
					wantItems[v.getName()] = append(wantItems[v.getName()], it)
				}
			}

			for k, v := range wantItems {
				require.ElementsMatch(s.T(), v, gotItems[k])
			}
		})
	}
}

func (s *suiteTest) Test_initGlobalHoarder() {
	tests := []struct {
		name               string
		initHoarders       []Hoarder
		expectedHoarderIdx int
	}{
		{
			name: "should be able to initialize a global hoarder",
			initHoarders: []Hoarder{
				&hoarder{
					inventoryMap: make(map[string]Inventory),
					mu:           sync.RWMutex{},
				},
			},
			expectedHoarderIdx: 0,
		},
		{
			name: "multiple initializations should not change the global hoarder",
			initHoarders: []Hoarder{
				&hoarder{
					inventoryMap: map[string]Inventory{
						defaultInventoryName: newInventory(defaultInventoryName),
					},
					mu: sync.RWMutex{},
				},
				&hoarder{
					inventoryMap: make(map[string]Inventory),
					mu:           sync.RWMutex{},
				},
			},
			expectedHoarderIdx: 0,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			once = sync.Once{}
			for _, h := range tt.initHoarders {
				initGlobalHoarder(h)
			}

			require.Same(s.T(), tt.initHoarders[tt.expectedHoarderIdx], globalHoarder)
		})
	}
}

func (s *suiteTest) Test_globalFactory() {
	tests := []struct {
		name              string
		numberOfInitCalls int
	}{
		{
			name:              "multiple calls should not change the global hoarder",
			numberOfInitCalls: 5,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			once = sync.Once{}
			globalHoarder = nil
			var got Hoarder
			for i := 0; i < tt.numberOfInitCalls; i++ {
				got = globalFactory()

				require.NotNil(s.T(), got)

				require.Same(s.T(), got, globalHoarder)
			}
		})
	}
}

func (s *suiteTest) Test_merge() {
	tests := []struct {
		name  string
		given []Hoarder
		want  Hoarder
	}{
		{
			name: "should be able to merge multiple hoarders",
			given: []Hoarder{
				&hoarder{
					mu: sync.RWMutex{},
				},
				&hoarder{
					inventoryMap: map[string]Inventory{
						"test": newInventory("test"),
						"test2": func() Inventory {
							i := newInventory("test2")
							i.Put(newItem("test thing11", "test11"))
							i.Put(newItem("test thing", "test"))

							return i
						}(),
					},
					mu: sync.RWMutex{},
				},
				&hoarder{
					inventoryMap: map[string]Inventory{
						"test": newInventory("test"),
						"test2": func() Inventory {
							i := newInventory("test2")
							i.Put(newItem("test thing234", "test234"))

							return i
						}(),
					},
					mu: sync.RWMutex{},
				},
				&hoarder{
					inventoryMap: map[string]Inventory{
						"test": newInventory("test"),
						"test2": func() Inventory {
							i := newInventory("test2")
							i.Put(newItem("test thing234", "test234"))

							return i
						}(),
					},
					mu: sync.RWMutex{},
				},
				&hoarder{
					inventoryMap: map[string]Inventory{
						"test53": newInventory("test53"),
						"test265": func() Inventory {
							i := newInventory("test265")
							i.Put(newItem("test thing23124", "test23124"))

							return i
						}(),
					},
					mu: sync.RWMutex{},
				},
			},
			want: &hoarder{
				inventoryMap: map[string]Inventory{
					"test": newInventory("test"),
					"test2": func() Inventory {
						i := newInventory("test2")
						i.Put(newItem("test thing", "test"))
						i.Put(newItem("test thing11", "test11"))
						i.Put(newItem("test thing234", "test234"))

						return i
					}(),
					"test53": newInventory("test53"),
					"test265": func() Inventory {
						i := newInventory("test265")
						i.Put(newItem("test thing23124", "test23124"))

						return i
					}(),
				},
				mu: sync.RWMutex{},
			},
		},
		{
			name: "should be able to handle merge with nil hoarder by returning the same hoarder",
			given: []Hoarder{
				&hoarder{
					inventoryMap: map[string]Inventory{
						"test": newInventory("test"),
						"test2": func() Inventory {
							i := newInventory("test2")
							i.Put(newItem("test thing234", "test234"))

							return i
						}(),
					},
					mu: sync.RWMutex{},
				},
				nil,
			},
			want: &hoarder{
				inventoryMap: map[string]Inventory{
					"test": newInventory("test"),
					"test2": func() Inventory {
						i := newInventory("test2")
						i.Put(newItem("test thing234", "test234"))

						return i
					}(),
				},
				mu: sync.RWMutex{},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			var h Hoarder

			if len(tt.given) > 0 {
				h = tt.given[0]
			} else {
				once = sync.Once{}
				initGlobalHoarder(factory())

				h = globalHoarder
			}

			for _, v := range tt.given {
				h.merge(v)
			}

			gotInventories := []Inventory{}
			gotInventoriesNames := []string{}
			for _, v := range h.loadout() {
				gotInventories = append(gotInventories, v)
				gotInventoriesNames = append(gotInventoriesNames, v.getName())
			}

			wantInventories := []Inventory{}
			wantInventoriesNames := []string{}
			for _, v := range tt.want.loadout() {
				wantInventories = append(wantInventories, v)
				wantInventoriesNames = append(wantInventoriesNames, v.getName())
			}

			require.ElementsMatch(s.T(), wantInventoriesNames, gotInventoriesNames)

			gotItems := map[string][]Item{}
			for _, v := range gotInventories {
				for _, it := range v.loadout() {
					if gotItems[v.getName()] == nil {
						gotItems[v.getName()] = []Item{}
					}
					gotItems[v.getName()] = append(gotItems[v.getName()], it)
				}
			}

			wantItems := map[string][]Item{}
			for _, v := range wantInventories {
				for _, it := range v.loadout() {
					if wantItems[v.getName()] == nil {
						wantItems[v.getName()] = []Item{}
					}
					wantItems[v.getName()] = append(wantItems[v.getName()], it)
				}
			}

			for k, v := range wantItems {
				require.ElementsMatch(s.T(), v, gotItems[k])
			}
		})
	}

}

func (s *suiteTest) Test_get() {
	tests := []struct {
		name               string
		givenHoarder       Hoarder
		givenType          reflect.Type
		givenInventoryName string
		givenItemName      string
		want               interface{}
	}{
		{
			name: "should be able to get a struct item from default inventory",
			givenHoarder: func() Hoarder {
				h := Hoard(nil, TestFooImpl{}, "test", "test2", "test3", &TestFooImpl{}, &hoarder{})

				return h
			}(),
			givenType:          reflect.TypeOf(TestFooImpl{}),
			givenInventoryName: defaultInventoryName,
			givenItemName:      "",
			want:               TestFooImpl{},
		},
		{
			name: "should be able to get a named struct item from default inventory",
			givenHoarder: func() Hoarder {
				h := Hoard(nil, TestFooImpl{}, "test", "test2", "test3", &TestFooImpl{}, &hoarder{}, RememberAs(TestFooImpl{}, "test"))

				return h
			}(),
			givenType:          reflect.TypeOf(TestFooImpl{}),
			givenInventoryName: defaultInventoryName,
			givenItemName:      "test",
			want:               TestFooImpl{},
		},
		{
			name: "should be able to get the latest item of the same type if multiple items of the same type hoarded",
			givenHoarder: func() Hoarder {
				h := Hoard(nil, TestFooImpl{}, "test", "test2", "test3", &TestFooImpl{}, &hoarder{}, RememberAs(TestFooImpl{}, "test"))

				return h
			}(),
			givenType:          reflect.TypeOf("TestFooImpl{}"),
			givenInventoryName: defaultInventoryName,
			givenItemName:      "",
			want:               "test3",
		},
		{
			name: "should be able to get implemented interface item from inventory",
			givenHoarder: func() Hoarder {
				h := Hoard(nil, TestFooImpl{}, "test", "test2", "test3", &hoarder{}, RememberAs(TestFooImpl{}, "test"))

				return h
			}(),
			givenType:          reflect.TypeFor[TestFooer](),
			givenInventoryName: defaultInventoryName,
			givenItemName:      "",
			want:               TestFooImpl{},
		},
		{
			name: "should be able to get item from named inventory within the default inventory if the item is not yet hoarded in the default inventory",
			givenHoarder: func() Hoarder {
				invent := UseInventory("test")
				invent.Put(RememberAs(TestFooImpl{}, "test"))
				h := factory(&hoarder{}, invent)

				return h
			}(),
			givenType:          reflect.TypeOf(TestFooImpl{}),
			givenInventoryName: defaultInventoryName,
			givenItemName:      "test",
			want:               TestFooImpl{},
		},
		{
			name: "should return nil if the item is not hoarded",
			givenHoarder: func() Hoarder {
				h := Hoard(nil, TestFooImpl{}, "test", "test2", "test3", &TestFooImpl{}, &hoarder{}, RememberAs(TestFooImpl{}, "test"))

				return h
			}(),
			givenType:          reflect.TypeOf(TestFooImpl{}),
			givenInventoryName: defaultInventoryName,
			givenItemName:      "test4",
			want:               nil,
		},
		{
			name: "should return nil if the requested item is a function",
			givenHoarder: func() Hoarder {
				invent := UseInventory("test")
				invent.Put(RememberAs(TestFooImpl{}, "test"))
				h := factory(&hoarder{}, invent)

				return h
			}(),
			givenType:          reflect.TypeOf(func() {}),
			givenInventoryName: defaultInventoryName,
			givenItemName:      "",
			want:               nil,
		},
		{
			name:         "should return nil if the requested inventory is not mapped",
			givenHoarder: Hoard(nil, &hoarder{}),
			givenType:    reflect.TypeOf(TestFooImpl{}),
			want:         nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got := tt.givenHoarder.get(tt.givenType, tt.givenInventoryName, tt.givenItemName)

			require.Equal(s.T(), tt.want, got)
		})
	}
}

func (s *suiteTest) TestEquipWithOption() {

	customHoarder := Hoard(
		HoardOptions{}.ShouldReplaceGlobal(false),
		UseInventory("test").
			Put(RememberAs(TestFooImpl{}, "test")).
			Put(RememberAs(TestFooImpl{}, "test2")).
			Put(RememberAs(TestFooImpl{}, "test3")).
			Put(RememberAs("test353", "")).
			Put(RememberAs(&TestFooImpl{}, "")),
	)

	tests := []struct {
		name                string
		givenOption         EquipOptions
		givenCustomHoarders []Hoarder
		wantString          string
		wantStruct          TestFooImpl
		wantPointer         *TestFooImpl
	}{
		{
			name:                "should be able to equip with no option with no custom hoarders",
			givenOption:         nil,
			givenCustomHoarders: []Hoarder{},
			wantString:          "test",
			wantStruct:          TestFooImpl{},
			wantPointer:         &TestFooImpl{},
		},
		{
			name:        "should be able to equip with options and custom hoarders",
			givenOption: EquipOptions{}.WithCustomInventoryName("test"),
			givenCustomHoarders: []Hoarder{
				Hoard(
					HoardOptions{}.ShouldReplaceGlobal(false).WithCustomHoarder(customHoarder),
					UseInventory("test").
						Put(RememberAs(TestFooImpl{}, "test")).
						Put(RememberAs(TestFooImpl{}, "test2")).
						Put(RememberAs(TestFooImpl{}, "test3")).
						Put(RememberAs("test353", "")).
						Put(RememberAs(&TestFooImpl{}, "")),
				),
			},
			wantString:  "test353",
			wantStruct:  TestFooImpl{},
			wantPointer: &TestFooImpl{},
		},
		{
			name:        "should be able to equip with options and custom hoarders and custom item name",
			givenOption: EquipOptions{}.WithCustomInventoryName("test").WithCustomItemName("test"),
			givenCustomHoarders: []Hoarder{
				Hoard(
					HoardOptions{}.ShouldReplaceGlobal(false),
					UseInventory("test").
						Put(RememberAs(TestFooImpl{}, "test")).
						Put(RememberAs(TestFooImpl{}, "test2")).
						Put(RememberAs(TestFooImpl{}, "test3")).
						Put(RememberAs("test353", "test")).
						Put(RememberAs(&TestFooImpl{}, "test")),
				),
			},
			wantString:  "test353",
			wantStruct:  TestFooImpl{},
			wantPointer: &TestFooImpl{},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			Hoard(nil, "test", TestFooImpl{}, &TestFooImpl{}, TestFooImpl{}, &hoarder{}, RememberAs(TestFooImpl{}, "test"))

			gotString := EquipWithOption[string](tt.givenOption, tt.givenCustomHoarders...)
			gotStruct := EquipWithOption[TestFooImpl](tt.givenOption, tt.givenCustomHoarders...)
			gotPointer := EquipWithOption[*TestFooImpl](tt.givenOption, tt.givenCustomHoarders...)

			require.Equal(s.T(), tt.wantString, gotString)
			require.Equal(s.T(), tt.wantStruct, gotStruct)
			require.Equal(s.T(), tt.wantPointer, gotPointer)
		})
	}
}

func (s *suiteTest) TestEquipDefault() {
	tests := []struct {
		name                string
		givenCustomHoarders []Hoarder
		wantString          string
		wantStruct          TestFooImpl
		wantPointer         *TestFooImpl
	}{
		{
			name:                "should be able to equip with no custom hoarders",
			givenCustomHoarders: []Hoarder{},
			wantString:          "test",
			wantStruct:          TestFooImpl{},
			wantPointer:         &TestFooImpl{},
		},
		{
			name: "should be able to equip with custom hoarders",
			givenCustomHoarders: []Hoarder{
				Hoard(
					HoardOptions{}.ShouldReplaceGlobal(false),
					UseInventory("test").
						Put(RememberAs(TestFooImpl{}, "test")).
						Put(RememberAs(TestFooImpl{}, "test2")).
						Put(RememberAs(TestFooImpl{}, "test3")).
						Put(RememberAs("test353", "")).
						Put(RememberAs(&TestFooImpl{}, "")),
				),
			},
			wantString:  "test353",
			wantStruct:  TestFooImpl{},
			wantPointer: &TestFooImpl{},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			Hoard(nil, "test", TestFooImpl{}, &TestFooImpl{}, TestFooImpl{}, &hoarder{}, RememberAs(TestFooImpl{}, "test"))

			gotString := EquipDefault[string](tt.givenCustomHoarders...)
			gotStruct := EquipDefault[TestFooImpl](tt.givenCustomHoarders...)
			gotPointer := EquipDefault[*TestFooImpl](tt.givenCustomHoarders...)

			require.Equal(s.T(), tt.wantString, gotString)
			require.Equal(s.T(), tt.wantStruct, gotStruct)
			require.Equal(s.T(), tt.wantPointer, gotPointer)
		})
	}
}

// no-op test to cover the loadout function's break statement
func (s *suiteTest) Test_loadout_break() {
	h := Hoard(nil, "test", TestFooImpl{}, &TestFooImpl{}, TestFooImpl{}, &hoarder{}, RememberAs(TestFooImpl{}, "test"))

	for range h.loadout() {
		break
	}
}
