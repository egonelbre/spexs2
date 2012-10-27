package utils

import (
	"reflect"
	"runtime"
)

func FuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
