package hoard

import (
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/oopchi/hoard/internal/core/constant"
	"github.com/oopchi/hoard/internal/domain/entity"
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
				inventoryMap: map[string]entity.Inventory{
					constant.DEFAULT_INVENTORY_NAME: entity.NewInventory(constant.DEFAULT_INVENTORY_NAME),
				},
				mu: sync.RWMutex{},
			},
		},
		{
			name: "should be able to create a hoarder with a single named item",
			given: []interface{}{
				entity.NewItem("test thing", "test"),
			},
			want: &hoarder{
				inventoryMap: func() map[string]entity.Inventory {
					inventoryMap := make(map[string]entity.Inventory)
					inventoryMap[constant.DEFAULT_INVENTORY_NAME] = entity.NewInventory(constant.DEFAULT_INVENTORY_NAME)
					inventoryMap[constant.DEFAULT_INVENTORY_NAME].Put(entity.NewItem("test thing", "test"))

					return inventoryMap
				}(),
				mu: sync.RWMutex{},
			},
		},
		{
			name: "should be able to create a hoarder with a single named inventory",
			given: []interface{}{
				entity.NewInventory("test"),
			},
			want: &hoarder{
				inventoryMap: func() map[string]entity.Inventory {
					inventoryMap := make(map[string]entity.Inventory)
					inventoryMap["test"] = entity.NewInventory("test")
					inventoryMap[constant.DEFAULT_INVENTORY_NAME] = entity.NewInventory(constant.DEFAULT_INVENTORY_NAME)

					return inventoryMap
				}(),
				mu: sync.RWMutex{},
			},
		},
		{
			name: "should be able to create a hoarder with a single named inventory with a single pointer item and a single named item",
			given: []interface{}{
				func() entity.Inventory {
					i := entity.NewInventory("test")
					i.Put(entity.NewItem(&TestFooImpl{}, "test234"))

					return i
				}(),
				entity.NewItem("test thing", "test"),
			},
			want: &hoarder{
				inventoryMap: func() map[string]entity.Inventory {
					inventoryMap := make(map[string]entity.Inventory)
					inventoryMap["test"] = entity.NewInventory("test")
					inventoryMap["test"].Put(entity.NewItem(&TestFooImpl{}, "test234"))
					inventoryMap[constant.DEFAULT_INVENTORY_NAME] = entity.NewInventory(constant.DEFAULT_INVENTORY_NAME).
						Put(entity.NewItem("test thing", "test")).
						Put(entity.NewItem(&TestFooImpl{}, "test234"))

					return inventoryMap
				}(),
				mu: sync.RWMutex{},
			},
		},
		{
			name: "should be able to create a hoarder with a single named inventory with a single pointer item and a single named item and a single interface item and a nil item and a struct item and a pointer item and a hoarder itself",
			given: []interface{}{
				func() entity.Inventory {
					i := entity.NewInventory("test")
					i.Put(entity.NewItem(&TestFooImpl{}, "test234"))

					return i
				}(),
				entity.NewItem("test thing", "test"),
				func() TestFooer {
					return TestFooImpl{}
				},
				"hehe",
				nil,
				TestFooImpl{},
				&TestFooImpl{},
				&hoarder{
					inventoryMap: map[string]entity.Inventory{
						"test": entity.NewInventory("test"),
						"test2": func() entity.Inventory {
							i := entity.NewInventory("test2")
							i.Put(entity.NewItem("test thing234", "test234"))

							return i
						}(),
					},
					mu: sync.RWMutex{},
				},
			},
			want: &hoarder{
				inventoryMap: func() map[string]entity.Inventory {
					inventoryMap := make(map[string]entity.Inventory)
					inventoryMap["test"] = entity.NewInventory("test")
					inventoryMap["test"].Put(entity.NewItem(&TestFooImpl{}, "test234"))
					inventoryMap[constant.DEFAULT_INVENTORY_NAME] = entity.NewInventory(constant.DEFAULT_INVENTORY_NAME)
					inventoryMap[constant.DEFAULT_INVENTORY_NAME].Put(entity.NewItem("test thing", "test"))
					inventoryMap[constant.DEFAULT_INVENTORY_NAME].Put(entity.NewItem(TestFooImpl{}, "github.com/oopchi/hoardTestFooImpl"))
					inventoryMap[constant.DEFAULT_INVENTORY_NAME].Put(entity.NewItem(&TestFooImpl{}, "*github.com/oopchi/hoardTestFooImpl"))
					inventoryMap[constant.DEFAULT_INVENTORY_NAME].Put(entity.NewItem("hehe", "string"))
					inventoryMap[constant.DEFAULT_INVENTORY_NAME].Put(entity.NewItem(&hoarder{
						inventoryMap: map[string]entity.Inventory{
							"test": entity.NewInventory("test"),
							"test2": func() entity.Inventory {
								i := entity.NewInventory("test2")
								i.Put(entity.NewItem("test thing234", "test234"))

								return i
							}(),
						},
						mu: sync.RWMutex{},
					}, "*github.com/oopchi/hoardhoarder")).Put(entity.NewItem(&TestFooImpl{}, "test234"))

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

			gotInventories := []entity.Inventory{}
			gotInventoriesNames := []string{}
			for _, v := range got.loadout() {
				gotInventories = append(gotInventories, v)
				gotInventoriesNames = append(gotInventoriesNames, v.GetName())
			}

			wantInventories := []entity.Inventory{}
			wantInventoriesNames := []string{}
			for _, v := range tt.want.loadout() {
				wantInventories = append(wantInventories, v)
				wantInventoriesNames = append(wantInventoriesNames, v.GetName())
			}

			require.ElementsMatch(s.T(), wantInventoriesNames, gotInventoriesNames)

			gotItems := map[string][]entity.Item{}
			for _, v := range gotInventories {
				for _, item := range v.Loadout() {
					if gotItems[v.GetName()] == nil {
						gotItems[v.GetName()] = []entity.Item{}
					}
					gotItems[v.GetName()] = append(gotItems[v.GetName()], item)
				}
			}

			wantItems := map[string][]entity.Item{}
			for _, v := range wantInventories {
				for _, item := range v.Loadout() {
					if wantItems[v.GetName()] == nil {
						wantItems[v.GetName()] = []entity.Item{}
					}
					wantItems[v.GetName()] = append(wantItems[v.GetName()], item)
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
					inventoryMap: make(map[string]entity.Inventory),
					mu:           sync.RWMutex{},
				},
			},
			expectedHoarderIdx: 0,
		},
		{
			name: "multiple initializations should not change the global hoarder",
			initHoarders: []Hoarder{
				&hoarder{
					inventoryMap: map[string]entity.Inventory{
						constant.DEFAULT_INVENTORY_NAME: entity.NewInventory(constant.DEFAULT_INVENTORY_NAME),
					},
					mu: sync.RWMutex{},
				},
				&hoarder{
					inventoryMap: make(map[string]entity.Inventory),
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
					inventoryMap: map[string]entity.Inventory{
						"test": entity.NewInventory("test"),
						"test2": func() entity.Inventory {
							i := entity.NewInventory("test2")
							i.Put(entity.NewItem("test thing11", "test11"))
							i.Put(entity.NewItem("test thing", "test"))

							return i
						}(),
					},
					mu: sync.RWMutex{},
				},
				&hoarder{
					inventoryMap: map[string]entity.Inventory{
						"test": entity.NewInventory("test"),
						"test2": func() entity.Inventory {
							i := entity.NewInventory("test2")
							i.Put(entity.NewItem("test thing234", "test234"))

							return i
						}(),
					},
					mu: sync.RWMutex{},
				},
				&hoarder{
					inventoryMap: map[string]entity.Inventory{
						"test": entity.NewInventory("test"),
						"test2": func() entity.Inventory {
							i := entity.NewInventory("test2")
							i.Put(entity.NewItem("test thing234", "test234"))

							return i
						}(),
					},
					mu: sync.RWMutex{},
				},
				&hoarder{
					inventoryMap: map[string]entity.Inventory{
						"test53": entity.NewInventory("test53"),
						"test265": func() entity.Inventory {
							i := entity.NewInventory("test265")
							i.Put(entity.NewItem("test thing23124", "test23124"))

							return i
						}(),
					},
					mu: sync.RWMutex{},
				},
			},
			want: &hoarder{
				inventoryMap: map[string]entity.Inventory{
					"test": entity.NewInventory("test"),
					"test2": func() entity.Inventory {
						i := entity.NewInventory("test2")
						i.Put(entity.NewItem("test thing", "test"))
						i.Put(entity.NewItem("test thing11", "test11"))
						i.Put(entity.NewItem("test thing234", "test234"))

						return i
					}(),
					"test53": entity.NewInventory("test53"),
					"test265": func() entity.Inventory {
						i := entity.NewInventory("test265")
						i.Put(entity.NewItem("test thing23124", "test23124"))

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
					inventoryMap: map[string]entity.Inventory{
						"test": entity.NewInventory("test"),
						"test2": func() entity.Inventory {
							i := entity.NewInventory("test2")
							i.Put(entity.NewItem("test thing234", "test234"))

							return i
						}(),
					},
					mu: sync.RWMutex{},
				},
				nil,
			},
			want: &hoarder{
				inventoryMap: map[string]entity.Inventory{
					"test": entity.NewInventory("test"),
					"test2": func() entity.Inventory {
						i := entity.NewInventory("test2")
						i.Put(entity.NewItem("test thing234", "test234"))

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

			gotInventories := []entity.Inventory{}
			gotInventoriesNames := []string{}
			for _, v := range h.loadout() {
				gotInventories = append(gotInventories, v)
				gotInventoriesNames = append(gotInventoriesNames, v.GetName())
			}

			wantInventories := []entity.Inventory{}
			wantInventoriesNames := []string{}
			for _, v := range tt.want.loadout() {
				wantInventories = append(wantInventories, v)
				wantInventoriesNames = append(wantInventoriesNames, v.GetName())
			}

			require.ElementsMatch(s.T(), wantInventoriesNames, gotInventoriesNames)

			gotItems := map[string][]entity.Item{}
			for _, v := range gotInventories {
				for _, item := range v.Loadout() {
					if gotItems[v.GetName()] == nil {
						gotItems[v.GetName()] = []entity.Item{}
					}
					gotItems[v.GetName()] = append(gotItems[v.GetName()], item)
				}
			}

			wantItems := map[string][]entity.Item{}
			for _, v := range wantInventories {
				for _, item := range v.Loadout() {
					if wantItems[v.GetName()] == nil {
						wantItems[v.GetName()] = []entity.Item{}
					}
					wantItems[v.GetName()] = append(wantItems[v.GetName()], item)
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
			givenInventoryName: constant.DEFAULT_INVENTORY_NAME,
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
			givenInventoryName: constant.DEFAULT_INVENTORY_NAME,
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
			givenInventoryName: constant.DEFAULT_INVENTORY_NAME,
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
			givenInventoryName: constant.DEFAULT_INVENTORY_NAME,
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
			givenInventoryName: constant.DEFAULT_INVENTORY_NAME,
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
			givenInventoryName: constant.DEFAULT_INVENTORY_NAME,
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
			givenInventoryName: constant.DEFAULT_INVENTORY_NAME,
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

func BenchmarkSingleHoard(b *testing.B) {
	once = sync.Once{}
	globalHoarder = nil

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Hoard(nil, &TestFooImpl{})
	}
}

func Benchmark10Hoards(b *testing.B) {
	once = sync.Once{}
	globalHoarder = nil

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Hoard(
			nil,
			&TestFooImpl{},
			&TestFooImpl{},
			&TestFooImpl{},
			&TestFooImpl{},
			&TestFooImpl{},
			&TestFooImpl{},
			&TestFooImpl{},
			&TestFooImpl{},
			&TestFooImpl{},
			&TestFooImpl{},
		)
	}
}

func BenchmarkSingleHoardWithoutReplaceGlobal(b *testing.B) {
	once = sync.Once{}
	globalHoarder = nil

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Hoard(HoardOptions{}.ShouldReplaceGlobal(false), &TestFooImpl{})
	}
}

func Benchmark10HoardsWithoutReplaceGlobal(b *testing.B) {
	once = sync.Once{}
	globalHoarder = nil

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Hoard(
			HoardOptions{}.ShouldReplaceGlobal(false),
			&TestFooImpl{},
			&TestFooImpl{},
			&TestFooImpl{},
			&TestFooImpl{},
			&TestFooImpl{},
			&TestFooImpl{},
			&TestFooImpl{},
			&TestFooImpl{},
			&TestFooImpl{},
			&TestFooImpl{},
		)
	}
}

func BenchmarkEquipDefault(b *testing.B) {
	simulateHugeHoard()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EquipDefault[string]()
	}
}

func BenchmarkEquipWithOption(b *testing.B) {
	simulateHugeHoard()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EquipWithOption[string](EquipOptions{}.WithCustomInventoryName("test").WithCustomItemName("test"))
	}
}

func BenchmarkEquipInterfaceDefault(b *testing.B) {
	simulateHugeHoard()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EquipDefault[TestFooer]()
	}
}

func BenchmarkEquipInterfaceWithOption(b *testing.B) {
	simulateHugeHoard()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EquipWithOption[TestFooer](EquipOptions{}.WithCustomInventoryName("test").WithCustomItemName("test"))
	}
}

func simulateHugeHoard() {
	once = sync.Once{}
	globalHoarder = nil

	for i := 0; i < 1000; i++ {
		Hoard(nil, RememberAs(&hoarder{}, fmt.Sprintf("test%d", i)))
	}

	Hoard(
		nil,
		UseInventory("test").
			Put(RememberAs("test353", "test")).
			Put(RememberAs(TestFooImpl{}, "test")),
	)

	for i := 1000; i < 2000; i++ {
		Hoard(nil, RememberAs(&hoarder{}, fmt.Sprintf("test%d", i)))
	}
}
