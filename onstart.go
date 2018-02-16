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

	var invokes []interface{}
	invokes = append(invokes, invoke.Interface())

	return invokeOption(invokes)
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

	var invokes []interface{}
	invokes = append(invokes, invoke.Interface())

	return invokeOption(invokes)
}

func createInvokeType(funcs ...interface{}) reflect.Type {
	var in []reflect.Type
	var out []reflect.Type

	// append lifecycle as the first args to the invoke func
	lifecycle := reflect.TypeOf(func(lifecycle Lifecycle) {}).In(0)
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
