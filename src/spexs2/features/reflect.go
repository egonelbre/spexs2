package features

import (
	"errors"
	"reflect"
	. "spexs"
)

func functionAndType(fn interface{}) (v reflect.Value, t reflect.Type, ok bool) {
	v = reflect.ValueOf(fn)
	ok = v.Kind() == reflect.Func
	if !ok {
		return
	}
	t = v.Type()
	return
}

func CallCreateWithArgs(function CreateFunc, args [][]int) (Feature, error) {
	fn, fnType, ok := functionAndType(function)
	if !ok {
		return nil, errors.New("argument is not a function")
	}

	if fnType.NumIn() != len(args) {
		return nil, errors.New("invalid number of arguments")
	}

	arguments := make([]reflect.Value, fnType.NumIn())
	for i := range args {
		arguments[i] = reflect.ValueOf(args[i])
	}
	result := fn.Call(arguments)
	inter := result[0].Interface()
	return inter.(Feature), nil
}
