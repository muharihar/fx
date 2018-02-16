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
	"reflect"
)

// OnStart makes it easy to execute funcs w/ already started types
func OnStart(funcs ...interface{}) Option {
	invokeType := createInvokeType(funcs...)
	invoke := reflect.MakeFunc(invokeType, func(args []reflect.Value) []reflect.Value {
		lifecycle := args[0].Interface().(Lifecycle)
		lifecycle.Append(Hook{
			OnStart: func(ctx context.Context) error {
				callInvokeFuncs(args, funcs...)
				return nil
			},
		})
		return []reflect.Value{}
	})
	return invokeOption([]interface{}{invoke.Interface()})
}

// OnStop makes it easy to hook into process shutdown
func OnStop(funcs ...interface{}) Option {
	invokeType := createInvokeType(funcs...)
	invoke := reflect.MakeFunc(invokeType, func(args []reflect.Value) []reflect.Value {
		lifecycle := args[0].Interface().(Lifecycle)
		lifecycle.Append(Hook{
			OnStop: func(ctx context.Context) error {
				callInvokeFuncs(args, funcs...)
				return nil
			},
		})
		return []reflect.Value{}
	})
	return invokeOption([]interface{}{invoke.Interface()})
}

func createInvokeType(funcs ...interface{}) reflect.Type {
	var in []reflect.Type
	var out []reflect.Type

	// append lifecycle as the first args to the invoke func
	lifecycle := reflect.TypeOf((*Lifecycle)(nil)).Elem()
	in = append(in, lifecycle)

	// append args [1:] using the args of all funcs passed
	for _, fn := range funcs {
		t := reflect.TypeOf(fn)
		for i := 0; i < t.NumIn(); i++ {
			in = append(in, t.In(i))
		}
	}

	return reflect.FuncOf(in, out, false)
}

func callInvokeFuncs(args []reflect.Value, funcs ...interface{}) {
	argFrom := 1
	argTo := 0

	for _, fn := range funcs {
		f := reflect.ValueOf(fn)

		// gather args and shift arg markers
		argTo += f.Type().NumIn()
		in := args[argFrom : argTo+1]
		argFrom += argTo

		f.Call(in)
	}
}
