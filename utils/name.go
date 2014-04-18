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

func UnqualifiedNameOf(obj interface{}) string {
	fqn := reflect.TypeOf(obj).String()
	// remove the package prefix
	s := path.Ext(fqn)[1:]
	return s
}
