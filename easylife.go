package fx

import (
	"context"
	"fmt"
	"reflect"
)

// OnStart makes it easy to execute funcs w/ already started types
func OnStart(funcs ...interface{}) Option {
	// build args to pass to invoke func
	var in []reflect.Type

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

	// create a func type using all "in" args
	var out []reflect.Type
	invokeType := reflect.FuncOf(in, out, false)

	// invoke func implementation
	invoke := reflect.MakeFunc(invokeType, func(args []reflect.Value) []reflect.Value {

		// extract lifecycle from arg 0
		lifecycle := args[0].Interface().(Lifecycle)

		// on lifecycle start
		lifecycle.Append(Hook{
			OnStart: func(ctx context.Context) error {

				// call all funcs with their args
				fmt.Println("CALLING FUNCS AT START")
				for _, fn := range funcs {
					fmt.Println("func:", fn)
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

// go test -run TestOnStart . -v
