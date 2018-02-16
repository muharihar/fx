package fx

import (
	"context"
	"reflect"
)

// OnStart makes it easy to execute funcs w/ already started types
func OnStart(funcs ...interface{}) Option {
	invokeType := createInvokeType(funcs...)
	invoke := reflect.MakeFunc(invokeType, func(args []reflect.Value) []reflect.Value {

		// extract lifecycle from arg 0
		lifecycle := args[0].Interface().(Lifecycle)

		// on lifecycle start
		lifecycle.Append(Hook{
			OnStart: func(ctx context.Context) error {

				// call all funcs with their args
				argBeg := 1
				argEnd := 0
				for _, fn := range funcs {
					f := reflect.ValueOf(fn)

					// gather args and shift argBeg and argEnd
					argEnd += f.Type().NumIn()
					in := args[argBeg : argEnd+1]
					argBeg += argEnd

					f.Call(in)
				}
				return nil
			},
		})

		var ret []reflect.Value
		return ret
	})

	// pass our 1 onstart invoke to invokeOption
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
