package hoard

import (
	"time"

	"github.com/stretchr/testify/require"
)

func (s *suiteTest) TestGetName() {
	require.Equal(s.T(), "test", s.invent.getName())
}

func (s *suiteTest) TestEquip() {
	tests := []struct {
		name      string
		putItem   item
		equipName string
		wantThing interface{}
		wantName  string
	}{
		{
			name:      "equipping with an existing item name should return the correct item",
			putItem:   newItem("test thing", "test"),
			equipName: "test",
			wantThing: "test thing",
			wantName:  "test",
		},
		{
			name:      "equipping non-existent item name in the map should return a nil item",
			putItem:   newItem("test thing", "test"),
			equipName: "wrong",
			wantThing: nil,
			wantName:  "",
		},
		{
			name: "should be able to equip structs",
			putItem: newItem(
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
			s.invent.put(tt.putItem)

			item := s.invent.equip(tt.equipName)

			if tt.wantThing == nil {
				require.Nil(s.T(), item)
				return
			}

			require.Equal(s.T(), tt.wantThing, item.use())

			require.NotNil(s.T(), item)
			require.Equal(s.T(), tt.wantName, item.getName())
		})
	}
}

func (s *suiteTest) TestPut() {
	tests := []struct {
		name      string
		putItem   item
		equipName string
		wantThing interface{}
		wantName  string
	}{
		{
			name:      "putting item should store the item in the inventory",
			putItem:   newItem("test thing", "test"),
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
			putItem: newItem(
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
			s.invent.put(tt.putItem)

			item := s.invent.equip(tt.equipName)

			if tt.wantThing == nil {
				require.Nil(s.T(), item)
				return
			}

			require.Equal(s.T(), tt.wantThing, item.use())

			require.NotNil(s.T(), item)
			require.Equal(s.T(), tt.wantName, item.getName())

		})
	}
}

func (s *suiteTest) TestPutIfAbsent() {
	tests := []struct {
		name      string
		putItem   item
		equipName string
		wantThing interface{}
		wantName  string
	}{
		{
			name:      "putting item that has not existed should store the item in the inventory",
			putItem:   newItem("test thing", "test"),
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
			putItem: newItem(
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
			putItem: newItem(
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
			s.invent.put(newItem("test thing23", "test23"))
			s.invent.putIfAbsent(tt.putItem)

			item := s.invent.equip(tt.equipName)

			if tt.wantThing == nil {
				require.Nil(s.T(), item)
				return
			}

			require.Equal(s.T(), tt.wantThing, item.use())

			require.NotNil(s.T(), item)
			require.Equal(s.T(), tt.wantName, item.getName())

		})
	}
}

func (s *suiteTest) TestMerge() {
	tests := []struct {
		name      string
		inventory inventory
		equipName string
		wantThing interface{}
		wantName  string
	}{
		{
			name: "merge with the same name inventory should succeed",
			inventory: func() inventory {
				i := newInventory("test")
				i.put(newItem("test thing", "test"))

				return i
			}(),
			equipName: "test",
			wantThing: "test thing",
			wantName:  "test",
		},
		{
			name: "merge with different name inventory should also merge the inventory",
			inventory: func() inventory {
				i := newInventory("test234")
				i.put(newItem("test thing", "test"))

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
			s.invent.put(newItem("test thing23", "test23"))
			s.invent.merge(tt.inventory)

			item := s.invent.equip(tt.equipName)

			if tt.wantThing == nil {
				require.Nil(s.T(), item)
				return
			}

			require.Equal(s.T(), tt.wantThing, item.use())

			require.NotNil(s.T(), item)
			require.Equal(s.T(), tt.wantName, item.getName())
		})
	}
}

func (s *suiteTest) TestLoadout() {
	tests := []struct {
		name          string
		itemsToPut    []item
		indexToBreak  int
		itemsToReturn []item
	}{
		{
			name: "should return all items in the inventory",
			itemsToPut: []item{
				newItem("test thing", "test"),
				newItem("test thing2", "test2"),
				newItem("test thing3", "test3"),
			},
			indexToBreak: -1,
			itemsToReturn: []item{
				newItem("test thing", "test"),
				newItem("test thing2", "test2"),
				newItem("test thing3", "test3"),
			},
		},
		{
			name: "should return all items in the inventory until the index to break",
			itemsToPut: []item{
				newItem("test thing", "test"),
				newItem("test thing2", "test2"),
				newItem("test thing3", "test3"),
			},
			indexToBreak: 1,
			itemsToReturn: []item{
				newItem("test thing", "test"),
				newItem("test thing2", "test2"),
			},
		},
		{
			name:          "should return empty items if the inventory is empty",
			itemsToPut:    []item{},
			indexToBreak:  -1,
			itemsToReturn: []item{},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			for _, item := range tt.itemsToPut {
				s.invent.put(item)
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
				s.invent.put(newItem("test thing4", "test4"))

				duration := time.Since(timeNow)
				durationChan <- duration
			}()

			var i int
			returnedItems := make([]item, 0)
			for _, k := range s.invent.loadout() {
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
