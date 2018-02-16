package fx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOnStart(t *testing.T) {
	app := New(
		Provide(newA, newB),
		OnStart(
			func(a a, b b) {
				assert.True(t, a.Started)
				assert.True(t, b.Started)
			},
			func(b b, a a) {
				assert.True(t, a.Started)
				assert.True(t, b.Started)
			},
		),
	)
	require.NoError(t, app.Err())
	err := app.Start(context.Background())
	require.NoError(t, err)
}

func TestOnStop(t *testing.T) {
	app := New(
		Provide(newA, newB),
		OnStop(
			func(a a, b b) {
				assert.True(t, a.Started)
				assert.True(t, b.Started)
			},
			func(b b, a a) {
				assert.True(t, a.Started)
				assert.True(t, b.Started)
			},
		),
	)
	require.NoError(t, app.Err())

	err := app.Start(context.Background())
	require.NoError(t, err)

	err = app.Stop(context.Background())
	require.NoError(t, err)
}

type a struct {
	Started bool
}

func newA(lifecycle Lifecycle) a {
	a := a{}
	lifecycle.Append(Hook{
		OnStart: func(ctx context.Context) error {
			a.Started = true
			return nil
		},
	})
	return a
}

type b struct {
	Started bool
}

func newB(lifecycle Lifecycle) b {
	b := b{}
	lifecycle.Append(Hook{
		OnStart: func(ctx context.Context) error {
			b.Started = true
			return nil
		},
	})
	return b
}
