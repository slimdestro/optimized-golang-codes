package main

import (
	"context"
	"fmt"
	"time"
)

// Saga is a type that represents a saga
type Saga struct {
	ID        string
	Steps     []Step
	Ctx       context.Context
	CancelCtx context.CancelFunc
}

// Step is a type that represents a step in a saga
type Step struct {
	Name string
	Func func(ctx context.Context) error
}

// NewSaga creates a new saga
func NewSaga(id string, steps []Step) *Saga {
	ctx, cancel := context.WithCancel(context.Background())
	return &Saga{
		ID:        id,
		Steps:     steps,
		Ctx:       ctx,
		CancelCtx: cancel,
	}
}

// Run runs the saga
func (s *Saga) Run() error {
	for _, step := range s.Steps {
		err := step.Func(s.Ctx)
		if err != nil {
			s.CancelCtx()
			return err
		}
	}
	return nil
}

// ExampleStep is an example step
func ExampleStep(ctx context.Context) error {
	fmt.Println("Example Step")
	time.Sleep(time.Second * 5)
	return nil
}

func main() {
	saga := NewSaga("example-saga", []Step{
		{
			Name: "Example Step",
			Func: ExampleStep,
		},
	})
	err := saga.Run()
	if err != nil {
		fmt.Println(err)
	}
}