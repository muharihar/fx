package fx

import (
	"context"
	"fmt"
	"reflect"
)

// OnStart makes it easy to execute funcs w/ already started types
func OnStart(funcs ...interface{}) Option {
	// starting w/ a func where lifecycle is the first argument
	lifecycle := reflect.ValueOf(func(lifecycle Lifecycle) {})

	// create a list of args from all the args in all the funcs passed
	var in []reflect.Type
	for _, fn := range funcs {
		t := reflect.TypeOf(fn)
		for i := 0; i < t.NumIn(); i++ {
			in = append(in, t.In(i))
		}
	}

	// create the func to be invoked by combining lifecycle and all args
	invoke := lifecycle

	// create a new func using combined, then call functionality during lifecycle start
	f := reflect.MakeFunc(invoke.Type(), func(args []reflect.Value) []reflect.Value {

		// extract lifecycle from arg 0
		lifecycle := args[0].Interface().(Lifecycle)
		lifecycle.Append(Hook{
			OnStart: func(ctx context.Context) error {

				// TODO for each fn of funcs, call w/ the correct args
				fmt.Println("HOLLLLLLLLLLLLLLERRRRRRRRRRRRRRRR")
				return nil
			},
		})
		invoke.Call(args)

		// TODO need to return anything at all?
		var ret []reflect.Value
		return ret
	})

	// replace funcs passed in w/ our lifecycle'd func
	var altered []interface{}
	altered = append(altered, f.Interface())

	return invokeOption(altered)
}
