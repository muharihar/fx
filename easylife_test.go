package fx

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOnStart(t *testing.T) {
	type A struct{}

	app := New(
		Provide(func() A { return A{} }),
		OnStart(func(a A) {
			fmt.Println("I got an", a)
		}),
	)

	require.NoError(t, app.Err())
	err := app.Start(context.Background())
	require.NoError(t, err)
}
