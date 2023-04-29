package model

import "fmt"

func IsSameType(a, b interface{}) bool {
	aType := fmt.Sprintf("%T", a)
	bType := fmt.Sprintf("%T", b)
	return aType == bType
}
