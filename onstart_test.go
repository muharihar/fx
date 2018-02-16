package fx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOnStartOnStop(t *testing.T) {
	called := 0
	app := New(
		Provide(newA, newB),
		OnStart(
			func(a *a, b *b) {
				assert.True(t, a.Started)
				assert.True(t, b.Started)
				called++
			},
			func(b *b, a *a) {
				assert.True(t, a.Started)
				assert.True(t, b.Started)
				called++
			},
		),
		OnStop(
			func(a *a, b *b) {
				called++
			},
			func(b *b, a *a) {
				called++
			},
		),
	)

	require.NoError(t, app.Err())
	require.NoError(t, app.Start(context.Background()))
	require.NoError(t, app.Stop(context.Background()))

	assert.Equal(t, 4, called)
}

type a struct {
	Started bool
}

func newA(lifecycle Lifecycle) *a {
	a := &a{}
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

func newB(lifecycle Lifecycle) *b {
	b := &b{}
	lifecycle.Append(Hook{
		OnStart: func(ctx context.Context) error {
			b.Started = true
			return nil
		},
	})
	return b
}
