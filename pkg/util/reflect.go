package util

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func GetFuncName(foo interface{}) string {
	fullFuncName := runtime.FuncForPC(reflect.ValueOf(foo).Pointer()).Name()
	parts := strings.Split(fullFuncName, ".")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func GetPackageNameOfFunc(foo interface{}) string {
	fullFuncName := runtime.FuncForPC(reflect.ValueOf(foo).Pointer()).Name()
	parts := strings.Split(fullFuncName, "/")
	if len(parts) == 0 {
		return ""
	}
	lastPart := parts[len(parts)-1]
	return strings.Split(lastPart, ".")[0]
}

func GetStructName(in interface{}) string {
	fullTypeName := fmt.Sprintf("%T", in)
	nameComponents := strings.Split(fullTypeName, ".")
	return nameComponents[len(nameComponents)-1]
}
