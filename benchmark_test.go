package hoard

import (
	"fmt"
	"sync"
	"testing"
)

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
