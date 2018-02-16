// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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
