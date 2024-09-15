package hoard_test

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

func ExampleHoard_withUberFX() {
	app := fx.New(
		fx.Provide(NewServiceA, NewServiceB),
		fx.Invoke(HoardServiceB),
	)

	app.Start(context.Background())

	fmt.Println("Equipped ServiceB with", hoard.EquipDefault[*ServiceB]().A.Name)
	// Output: Equipped ServiceB with Uber FX Service
}
