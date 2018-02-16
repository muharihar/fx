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

	app := New(
		Provide(func() A { return A{Name: "Grayson"} }),
		OnStart(func(a A) {
			fmt.Println("I got an", a)
		}),
	)

	require.NoError(t, app.Err())
	err := app.Start(context.Background())
	require.NoError(t, err)
}
