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

	// lifecycle is first argument
	in = append(in, lifecycle.Type().In(0))

	// args from passed funcs are args [1:]
	for _, fn := range funcs {
		t := reflect.TypeOf(fn)
		for i := 0; i < t.NumIn(); i++ {
			in = append(in, t.In(i))
		}
	}

	// create the func to be invoked by combining lifecycle and all args
	//invoke := lifecycle
	var out []reflect.Type
	invoke := reflect.FuncOf(in, out, false)

	// create a new func using combined, then call functionality during lifecycle start
	f := reflect.MakeFunc(invoke, func(args []reflect.Value) []reflect.Value {

		// extract lifecycle from arg 0
		lifecycle := args[0].Interface().(Lifecycle)

		// call all passed funcs on lifecycle start
		lifecycle.Append(Hook{
			OnStart: func(ctx context.Context) error {
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

	// replace funcs passed in w/ our lifecycle'd func
	var altered []interface{}
	altered = append(altered, f.Interface())

	return invokeOption(altered)
}
