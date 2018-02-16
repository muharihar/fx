package fx

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOnStart(t *testing.T) {
	type A struct {
		Name string
	}
	type B struct {
		Name string
	}

	app := New(
		Provide(func() A { return A{Name: "Grayson"} }),
		Provide(func() B { return B{Name: "Abhinav"} }),

		OnStart(
			func(a A, b B) {
				fmt.Println("I got an", a, b)
			},
			func(b B, a A) {
				fmt.Println("I got a", b, a)
			},
		),
	)

	require.NoError(t, app.Err())
	err := app.Start(context.Background())
	require.NoError(t, err)
}

func TestOnStop(t *testing.T) {
	type A struct {
		Name string
	}
	type B struct {
		Name string
	}

	app := New(
		Provide(func() A { return A{Name: "Grayson"} }),
		Provide(func() B { return B{Name: "Abhinav"} }),

		OnStop(
			func(a A, b B) {
				fmt.Println("I got an", a, b)
			},
			func(b B, a A) {
				fmt.Println("I got a", b, a)
			},
		),
	)

	require.NoError(t, app.Err())

	err := app.Start(context.Background())
	require.NoError(t, err)

	err = app.Stop(context.Background())
	require.NoError(t, err)
}
