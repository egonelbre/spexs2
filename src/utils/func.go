package utils

import (
	"path"
	"reflect"
	"runtime"
)

func FuncFullName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

func FuncName(fn interface{}) string {
	return path.Ext(FuncFullName(fn))[1:]
}
