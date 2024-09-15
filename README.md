
# Hoard

## Introduction

[![Go Reference](https://pkg.go.dev/badge/github.com/oopchi/hoard.svg)](https://pkg.go.dev/github.com/oopchi/hoard)
[![Go Report Card](https://goreportcard.com/badge/github.com/oopchi/hoard)](https://goreportcard.com/report/github.com/oopchi/hoard)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

**Hoard** is a simple and highly concurrent service container for Golang. It allows you to:

- "Hoard" (register) items
- "Equip" (retrieve) items

The package supports features like annotations ("remember as") for handling multiple items of the same type, and custom groupings ("use custom Inventory"). 

Hoard is optimized for concurrent access and works with various data types. However, it currently **does not** support hoarding or equipping functions.

### Motivation Behind Hoard

**Hoard** was developed with a specific motivation: to address some of the limitations found in existing dependency injection frameworks like **Uber FX**. While Uber FX is highly regarded for its robust dependency injection capabilities, it often requires the use of **constructors** to inject services into methods or function bodies, which can introduce unnecessary complexity.

**Hoard** takes a different approach by allowing direct interaction with objects, eliminating the need for constructors. You can register and retrieve items (services) directly within your code without having to set up specific constructors. The focus of **Hoard** is on simplifying the management of dependencies, leaving object construction outside the scope of the package.

However, **Hoard** is by no means a replacement for **Uber FX**. In fact, both can work hand-in-hand. You can continue using Uber FX's powerful dependency injection framework while leveraging **Hoard** for scenarios where direct object interaction is more efficient or desirable.

#### Example: Using Hoard and Uber FX Together

Here’s how you can combine **Hoard** and **Uber FX** in a project:

```go
package main

import (
	"context"
	"fmt"

	"github.com/oopchi/hoard"
	"go.uber.org/fx"
)

// ServiceA is an example service that we want to hoard and inject.
type ServiceA struct {
	Name string
}

// ServiceB uses ServiceA.
type ServiceB struct {
	A *ServiceA
}

// NewServiceA constructor for Uber FX.
func NewServiceA() *ServiceA {
	return &ServiceA{Name: "Uber FX Service"}
}

// NewServiceB constructor for Uber FX that depends on ServiceA.
func NewServiceB(a *ServiceA) *ServiceB {
	return &ServiceB{A: a}
}

// HoardServiceB hoards ServiceB for future use in the application.
func HoardServiceB(lc fx.Lifecycle, s *ServiceB) {
	// Hoarding ServiceB
	hoard.Hoard(nil, s)
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			fmt.Println("ServiceB stopping...")
			return nil
		},
	})
}

func main() {
	app := fx.New(
		fx.Provide(NewServiceA, NewServiceB),
		fx.Invoke(HoardServiceB),
	)

	app.Start(context.Background())

	fmt.Println("Equipped ServiceB with", hoard.EquipDefault[*ServiceB]().A.Name)
}
```

In this example:
- **Uber FX** is used for setting up constructors and lifecycle management of services.
- **Hoard** complements it by hoarding the service (`ServiceB`) for use later in the program.
- You can easily "equip" `ServiceB` whenever it's needed, blending both frameworks to get the best of both worlds.

This integration allows for the flexibility and powerful lifecycle management of Uber FX alongside the simplicity of Hoard’s direct service access.

### The Idea Behind Hoard

The concept of **Hoard** is inspired by the idea of a **hoarder** in an RPG game—someone who collects and stores various items in their Inventory. This analogy is used to describe how **Hoard** manages dependency injection and service management in a Golang application:

- **Hoarder**: Represents the entity that manages and organizes items (i.e., services). Just like an RPG hoarder collects and stores items, a **Hoarder** accumulates and maintains various services.
  
- **Hoard**: The act of adding items (services) to a hoarder’s Inventory. In technical terms, this corresponds to **service registration** or **dependency injection**.

- **Items**: These are the services or dependencies you register within the hoarder. Similar to how an RPG hoarder stores different items, the **Hoard** manages various dependencies.

- **Inventory**: Represents the collection of items managed by the hoarder. This is analogous to an RPG Inventory, where different items are grouped and accessed as needed. A hoarder can maintain multiple inventories, each serving a specific purpose or grouping.

- **Equip**: Refers to the process of retrieving items (services) from the hoarder’s Inventory. In technical terms, this is known as **service discovery** or **retrieval**, where you access the registered services.

This framework simplifies dependency injection and service management in Golang applications by providing a structured and efficient way to handle various services and dependencies in a concurrent environment.

## Features

- **Service Hoarding**: Register services, including basic data types like `int`, `string`, `struct`, `pointer`, `boolean`, etc.
- **Annotations**: Use annotations to differentiate services of the same type, allowing you to "remember as" unique names.
- **Custom Inventory**: Group services in different "inventories" to isolate retrieval contexts.
- **Global Hoarder Replacement**: The global hoarder can be replaced automatically, or this behavior can be disabled via options.
- **100% Test Coverage**: The package includes thorough tests, ensuring stability.
- **Optimized for Concurrency**: Hoard is designed for efficient, concurrent usage across multiple goroutines.

## Installation

To install Hoard, use `go get`:

```bash
go get github.com/oopchi/hoard@latest
```

Import the package into your Go code:

```go
import "github.com/oopchi/hoard"
```

### Requirements

**Hoard** requires Go version 1.23 or higher. This is because **Hoard** utilizes Go 1.23's new `range` over function feature. Ensure that your Go environment is up-to-date to take advantage of all features and improvements provided by **Hoard**.

## Usage

Below are some examples of how to use Hoard in your Go project to get you started:

### Basic Hoarding and Equipping

Hoard services of various types such as `int`, `string`, `struct`, `pointer`, and `boolean`. If multiple services of the same type are hoarded without annotation, the latest one will override the older one.

```go
package main

import (
	"fmt"

	"github.com/oopchi/hoard"
)

type MyService struct {
	Name string
}

func main() {
	// Hoard items of different types
	hoard.Hoard(nil, 42, "Hoarding a string", true, &MyService{Name: "Service 1"})

	// Equip the items
	myInt := hoard.EquipDefault[int]()
	myString := hoard.EquipDefault[string]()
	myBool := hoard.EquipDefault[bool]()
	myService := hoard.EquipDefault[*MyService]()

	fmt.Println(myInt)          // Output: 42
	fmt.Println(myString)       // Output: Hoarding a string
	fmt.Println(myBool)         // Output: true
	fmt.Println(myService.Name) // Output: Service 1
}
```

### Custom Hoarder Example

You can also create a custom hoarder and disable the automatic replacement of the global hoarder:

```go
package main

import (
	"fmt"

	"github.com/oopchi/hoard"
)

func main() {
	// Disable global hoarder replacement
	options := hoard.HoardOptions{}.ShouldReplaceGlobal(false)

	// Hoard services with a custom hoarder
	customHoarder := hoard.Hoard(options, 100, "custom string")

	// Equip services from the custom hoarder
	i := hoard.EquipWithOption[int](nil, customHoarder)
	s := hoard.EquipWithOption[string](nil, customHoarder)

	fmt.Println(i) // 100
	fmt.Println(s) // custom string
}
```

### Hoarding Multiple Data Types

Hoard supports a variety of types, including `int`, `string`, `struct`, `pointer`, `boolean`, and more. Here's an example of registering and equipping multiple types:

```go
package main

import (
	"fmt"

	"github.com/oopchi/hoard"
)

type MyStruct struct {
	Name string
}

func main() {
	// Hoard services of different types
	hoard.Hoard(nil, 123, "test", true, MyStruct{Name: "hoard"}, &MyStruct{Name: "pointer"})

	// Equip services
	i := hoard.EquipDefault[int]()
	s := hoard.EquipDefault[string]()
	b := hoard.EquipDefault[bool]()
	st := hoard.EquipDefault[MyStruct]()
	p := hoard.EquipDefault[*MyStruct]()

	fmt.Println(i)  // 123
	fmt.Println(s)  // test
	fmt.Println(b)  // true
	fmt.Println(st) // {hoard}
	fmt.Println(p)  // &{pointer}
}

```

### Handling Multiple Services of the Same Type

When registering multiple services of the same type without annotations, the newer service will override the older one. To avoid this, use annotations:

```go
package main

import (
	"fmt"

	"github.com/oopchi/hoard"
)

type Service struct {
	Name string
}

func main() {
	// Hoard multiple services of the same type
	hoard.Hoard(nil, hoard.RememberAs(Service{Name: "Service 1"}, "service1"))
	hoard.Hoard(nil, hoard.RememberAs(Service{Name: "Service 2"}, "service2"))

	// Equip services by name
	s1 := hoard.EquipWithOption[Service](hoard.EquipOptions{}.WithCustomItemName("service1"))
	s2 := hoard.EquipWithOption[Service](hoard.EquipOptions{}.WithCustomItemName("service2"))

	fmt.Println(s1.Name) // Service 1
	fmt.Println(s2.Name) // Service 2
}

```

### Equipping Interfaces with Annotations

When hoarding and equipping interfaces, it's recommended to use annotations for better performance and clarity. Make sure to annotate each interface with a unique name to avoid issues, as the underlying implementation type doesn't differentiate between them.

```go
package main

import (
    "fmt"
    "github.com/oopchi/hoard"
)

type Service interface {
    Execute()
}

type ServiceA struct{}
func (s ServiceA) Execute() { fmt.Println("Service A") }

type ServiceB struct{}
func (s ServiceB) Execute() { fmt.Println("Service B") }

func main() {
    // Hoard services and annotate them
    hoard.Hoard(nil, hoard.RememberAs(ServiceA{}, "serviceA"))
    hoard.Hoard(nil, hoard.RememberAs(ServiceB{}, "serviceB"))

    // Equip services by their annotation
    sA := hoard.EquipWithOption[Service](hoard.EquipOptions{}.WithCustomItemName("serviceA"))
    sB := hoard.EquipWithOption[Service](hoard.EquipOptions{}.WithCustomItemName("serviceB"))

    sA.Execute() // Service A
    sB.Execute() // Service B
}
```

### Using Annotations and Custom Inventory

For better performance when working with interfaces, it's recommended to annotate services with unique names since the underlying object type doesn't differentiate between implementations.

```go
package main

import (
	"fmt"

	"github.com/oopchi/hoard"
)

type MyService interface {
	Execute()
}

type ServiceImpl struct {
	ID int
}

func (s ServiceImpl) Execute() {
	fmt.Printf("Executing Service with ID: %d\n", s.ID)
}

func main() {
	// Use custom Inventory and annotation to hoard multiple implementations
	hoard.Hoard(nil,
		hoard.UseInventory("legendary items inventory").
			Put(hoard.RememberAs(ServiceImpl{ID: 1}, "impl1")).
			Put(hoard.RememberAs(ServiceImpl{ID: 1}, "")). // pass an empty string to use the default item name
			Put(hoard.RememberAs(ServiceImpl{ID: 2}, "impl2")),
	)

	// Equip the services using annotations
	svc1 := hoard.EquipWithOption[MyService](hoard.EquipOptions{}.WithCustomItemName("impl1").WithCustomInventoryName("legendary items inventory"))

	// You can also skip the inventory name because if an item is just hoarded for the first time (no other same item has been hoarded yet)
	// then no matter the custom inventory used, it will also be stored at the default inventory
	svc2 := hoard.EquipWithOption[MyService](hoard.EquipOptions{}.WithCustomItemName("impl2"))

	// You can even skip the annotation whatsoever if there has only ever been one such item being hoarded even if its annotated
	// If there were already multiple such items being hoarded though, if its being hoarded through a custom inventory
	// then it won't override the one at the default inventory anymore, however it will still override the one at that custom inventory if any existed
	svc3 := hoard.EquipDefault[ServiceImpl]()

	// You can however override the default inventory again if you specifically hoard on the default inventory (hoarding without specifying [hoard.UseInventory])
	hoard.Hoard(nil, ServiceImpl{ID: 5}, ServiceImpl{ID: 8})

	svc4 := hoard.EquipDefault[ServiceImpl]()

	svc1.Execute() // Output: Executing Service with ID: 1
	svc2.Execute() // Output: Executing Service with ID: 2
	svc3.Execute() // Output: Executing Service with ID: 1
	svc4.Execute() // Output: Executing Service with ID: 8
}
```

### Disabling Global Hoarder Replacement

You can disable automatic replacement of the global hoarder using `HoardOptions`.

```go
package main

import (
	"fmt"

	"github.com/oopchi/hoard"
)

func main() {
	// Hoard with the option to disable global hoarder replacement
	customHoarder := hoard.Hoard(hoard.HoardOptions{}.ShouldReplaceGlobal(false), 42)
	customHoarder = hoard.Hoard(hoard.HoardOptions{}.ShouldReplaceGlobal(false).WithCustomHoarder(customHoarder), 42)
	

	// Hoard with global hoarder replaced
	hoard.Hoard(nil, 50)
	hoard.Hoard(nil, 70)

	// Equip from the custom hoarder without affecting the global one
	myInt := hoard.EquipDefault[int](customHoarder)
	fmt.Println(myInt) // Output: 42

	// Equip from default hoarder will search for items within the global hoarder
	myInt = hoard.EquipDefault[int]()
	fmt.Println(myInt) // Output: 70
}
```

### Take Note: Panics on Non-Registered Items

When attempting to equip a service that hasn't been hoarded, the `Equip` function **may panic**. Ensure that the services you are trying to equip have been properly registered to avoid runtime errors.

```go
package main

import (
	"fmt"
	"github.com/oopchi/hoard"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	// Trying to equip a service that wasn't hoarded, this will panic
	hoard.EquipDefault[int]()
    // Output: Recovered from panic: interface conversion: interface {} is nil, not int
}
```

### Limitation

Hoard currently does not support hoarding or equipping **functions**.

## Benchmarks

The following table provides benchmark results comparing performance with different configurations. Using options, such as disabling global hoarder replacement, can significantly improve performance.

| Benchmark                                | Iterations | Time (ns/op) | Memory (B/op) | Allocations (allocs/op) |
|------------------------------------------|------------|--------------|---------------|--------------------------|
| `BenchmarkSingleHoard-12`                | 100        | 9983         | 1022          | 16                       |
| `Benchmark10Hoards-12`                   | 100        | 61028        | 1887          | 43                       |
| `BenchmarkSingleHoardWithoutReplaceGlobal-12` | 100  | 5556         | 944           | 14                       |
| `Benchmark10HoardsWithoutReplaceGlobal-12`    | 100  | 33288        | 1808          | 41                       |
| `BenchmarkEquipDefault-12`               | 100        | 3086         | 48            | 2                        |
| `BenchmarkEquipWithOption-12`            | 100        | 4833         | 168           | 10                       |
| `BenchmarkEquipInterfaceDefault-12`      | 100        | 1193682      | 192           | 8                        |
| `BenchmarkEquipInterfaceWithOption-12`   | 100        | 5955         | 296           | 13                       |

### Statistical Remarks

- **Disabling Global Hoarder Replacement**: Disabling the global hoarder replacement (`BenchmarkSingleHoardWithoutReplaceGlobal`) results in a **44% improvement** in execution time compared to the default behavior.
- **Multiple Hoards**: Hoarding multiple items (e.g., `Benchmark10Hoards`) incurs higher memory usage and execution time due to managing more services, but you can reduce overhead by disabling replace global option.
- **Equip Performance**: Using `EquipWithOption` is slightly slower than `EquipDefault`, but it provides flexibility in selecting specific services by annotations or custom inventories, this however doesn't apply when trying to equip interfaces.
- **Interface Equipping**: Equipping interfaces without annotations (`BenchmarkEquipInterfaceDefault`) is much slower due to reflection and lack of type differentiation. Using annotations (`BenchmarkEquipInterfaceWithOption`) improves performance drastically by **~199x**.

## Documentation

Complete documentation can be found on [pkg.go.dev](https://pkg.go.dev/github.com/oopchi/hoard).

## Contributing

Contributions are welcome! Please open an issue or submit a pull request if you find a bug or think of an improvement.

### Steps to contribute:

1. Fork the repository.
2. Create a new branch: `git checkout -b feature-branch`.
3. Make your changes.
4. Submit a pull request.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](./LICENSE) file for details.

## Author

**Calvin Alfredo**  
[![GitHub](https://img.shields.io/badge/GitHub-333?style=for-the-badge&logo=github&logoColor=white)](https://github.com/oopchi)  
[![LinkedIn](https://img.shields.io/badge/LinkedIn-0A66C2?style=for-the-badge&logo=linkedin&logoColor=white)](https://www.linkedin.com/in/calvinalfrido/)

## Acknowledgments

- Inspired by the flexibility of [Uber FX](https://github.com/uber-go/fx) and the Go programming language.
