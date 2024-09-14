package entity

import (
	"time"

	"github.com/stretchr/testify/require"
)

func (s *suiteTest) TestGetName() {
	require.Equal(s.T(), "test", s.inventory.GetName())
}

func (s *suiteTest) TestEquip() {
	tests := []struct {
		name      string
		putItem   Item
		equipName string
		wantThing interface{}
		wantName  string
	}{
		{
			name:      "equipping with an existing item name should return the correct item",
			putItem:   NewItem("test thing", "test"),
			equipName: "test",
			wantThing: "test thing",
			wantName:  "test",
		},
		{
			name:      "equipping non-existent item name in the map should return a nil item",
			putItem:   NewItem("test thing", "test"),
			equipName: "wrong",
			wantThing: nil,
			wantName:  "",
		},
		{
			name: "should be able to equip structs",
			putItem: NewItem(
				struct {
					name string
				}{
					name: "test thing",
				},
				"test",
			),
			equipName: "test",
			wantThing: struct {
				name string
			}{
				name: "test thing",
			},
			wantName: "test",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.inventory.Put(tt.putItem)

			item := s.inventory.Equip(tt.equipName)

			if tt.wantThing == nil {
				require.Nil(s.T(), item)
				return
			}

			require.Equal(s.T(), tt.wantThing, item.Use())

			require.NotNil(s.T(), item)
			require.Equal(s.T(), tt.wantName, item.GetName())
		})
	}
}

func (s *suiteTest) TestPut() {
	tests := []struct {
		name      string
		putItem   Item
		equipName string
		wantThing interface{}
		wantName  string
	}{
		{
			name:      "putting item should store the item in the inventory",
			putItem:   NewItem("test thing", "test"),
			equipName: "test",
			wantThing: "test thing",
			wantName:  "test",
		},
		{
			name:      "should handle putting nil item by returning nil",
			putItem:   nil,
			equipName: "test",
			wantThing: nil,
			wantName:  "",
		},
		{
			name: "should be able to put structs",
			putItem: NewItem(
				struct {
					name string
				}{
					name: "test thing",
				},
				"test",
			),
			equipName: "test",
			wantThing: struct {
				name string
			}{
				name: "test thing",
			},
			wantName: "test",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.inventory.Put(tt.putItem)

			item := s.inventory.Equip(tt.equipName)

			if tt.wantThing == nil {
				require.Nil(s.T(), item)
				return
			}

			require.Equal(s.T(), tt.wantThing, item.Use())

			require.NotNil(s.T(), item)
			require.Equal(s.T(), tt.wantName, item.GetName())

		})
	}
}

func (s *suiteTest) TestPutIfAbsent() {
	tests := []struct {
		name      string
		putItem   Item
		equipName string
		wantThing interface{}
		wantName  string
	}{
		{
			name:      "putting item that has not existed should store the item in the inventory",
			putItem:   NewItem("test thing", "test"),
			equipName: "test",
			wantThing: "test thing",
			wantName:  "test",
		},
		{
			name:      "should handle putting nil item by returning nil",
			putItem:   nil,
			equipName: "test",
			wantThing: nil,
			wantName:  "",
		},
		{
			name: "should be able to put structs",
			putItem: NewItem(
				struct {
					name string
				}{
					name: "test thing",
				},
				"test",
			),
			equipName: "test",
			wantThing: struct {
				name string
			}{
				name: "test thing",
			},
			wantName: "test",
		},
		{
			name: "should not store the item if the item already exists",
			putItem: NewItem(
				struct {
					name string
				}{
					name: "test thing",
				},
				"test23",
			),
			equipName: "test23",
			wantThing: "test thing23",
			wantName:  "test23",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.inventory.Put(NewItem("test thing23", "test23"))
			s.inventory.PutIfAbsent(tt.putItem)

			item := s.inventory.Equip(tt.equipName)

			if tt.wantThing == nil {
				require.Nil(s.T(), item)
				return
			}

			require.Equal(s.T(), tt.wantThing, item.Use())

			require.NotNil(s.T(), item)
			require.Equal(s.T(), tt.wantName, item.GetName())

		})
	}
}

func (s *suiteTest) TestMerge() {
	tests := []struct {
		name      string
		inventory Inventory
		equipName string
		wantThing interface{}
		wantName  string
	}{
		{
			name: "merge with the same name inventory should succeed",
			inventory: func() Inventory {
				i := NewInventory("test")
				i.Put(NewItem("test thing", "test"))

				return i
			}(),
			equipName: "test",
			wantThing: "test thing",
			wantName:  "test",
		},
		{
			name: "merge with different name inventory should also merge the inventory",
			inventory: func() Inventory {
				i := NewInventory("test234")
				i.Put(NewItem("test thing", "test"))

				return i
			}(),
			equipName: "test",
			wantThing: "test thing",
			wantName:  "test",
		},
		{
			name:      "merge with nil inventory should not affect the current inventory",
			inventory: nil,
			equipName: "test23",
			wantThing: "test thing23",
			wantName:  "test23",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.inventory.Put(NewItem("test thing23", "test23"))
			s.inventory.Merge(tt.inventory)

			item := s.inventory.Equip(tt.equipName)

			if tt.wantThing == nil {
				require.Nil(s.T(), item)
				return
			}

			require.Equal(s.T(), tt.wantThing, item.Use())

			require.NotNil(s.T(), item)
			require.Equal(s.T(), tt.wantName, item.GetName())
		})
	}
}

func (s *suiteTest) TestLoadout() {
	tests := []struct {
		name          string
		itemsToPut    []Item
		indexToBreak  int
		itemsToReturn []Item
	}{
		{
			name: "should return all items in the inventory",
			itemsToPut: []Item{
				NewItem("test thing", "test"),
				NewItem("test thing2", "test2"),
				NewItem("test thing3", "test3"),
			},
			indexToBreak: -1,
			itemsToReturn: []Item{
				NewItem("test thing", "test"),
				NewItem("test thing2", "test2"),
				NewItem("test thing3", "test3"),
			},
		},
		{
			name: "should return all items in the inventory until the index to break",
			itemsToPut: []Item{
				NewItem("test thing", "test"),
				NewItem("test thing2", "test2"),
				NewItem("test thing3", "test3"),
			},
			indexToBreak: 1,
			itemsToReturn: []Item{
				NewItem("test thing", "test"),
				NewItem("test thing2", "test2"),
			},
		},
		{
			name:          "should return empty items if the inventory is empty",
			itemsToPut:    []Item{},
			indexToBreak:  -1,
			itemsToReturn: []Item{},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			for _, item := range tt.itemsToPut {
				s.inventory.Put(item)
			}

			durationChan := make(chan time.Duration)
			defer close(durationChan)

			startChan := make(chan time.Time)
			isStartChanClosed := false
			defer func() {
				if !isStartChanClosed {
					close(startChan)
				}
			}()

			go func() {
				if len(tt.itemsToReturn) == 0 {
					durationChan <- 0
					return
				}
				timeNow := <-startChan
				s.inventory.Put(NewItem("test thing4", "test4"))

				duration := time.Since(timeNow)
				durationChan <- duration
			}()

			var i int
			returnedItems := make([]Item, 0)
			for _, k := range s.inventory.Loadout() {
				if i == 0 {
					startChan <- time.Now()
					isStartChanClosed = true
				}
				time.Sleep(1 * time.Second)
				returnedItems = append(returnedItems, k)
				if i == tt.indexToBreak {
					break
				}
				i++
			}

			require.Equal(s.T(), len(tt.itemsToReturn), len(returnedItems))
			duration := <-durationChan

			expectedDuration := time.Duration(len(returnedItems)) * time.Second
			require.GreaterOrEqual(s.T(), duration, expectedDuration)
		})
	}
}

func (s *suiteTest) TestClone() {
	tests := []struct {
		name      string
		inventory Inventory
		equipName string
		wantThing interface{}
		wantName  string
	}{
		{
			name: "clone an inventory should return a new inventory with the same items",
			inventory: func() Inventory {
				i := NewInventory("test")
				i.Put(NewItem("test thing", "test"))

				return i
			}(),
			equipName: "test",
			wantThing: "test thing",
			wantName:  "test",
		},
		{
			name:      "clone an empty inventory should return a new empty inventory",
			inventory: NewInventory("test"),
			equipName: "test",
			wantThing: nil,
			wantName:  "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.inventory.Put(NewItem("test thing23", "test23"))
			clonedInventory := tt.inventory.Clone()

			item := clonedInventory.Equip(tt.equipName)

			if tt.wantThing == nil {
				require.Nil(s.T(), item)
				return
			}

			require.Equal(s.T(), tt.wantThing, item.Use())

			require.NotNil(s.T(), item)
			require.Equal(s.T(), tt.wantName, item.GetName())
		})
	}
}
