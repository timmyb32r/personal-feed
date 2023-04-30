package util

import (
	"fmt"
	"strings"
)

func GetStructName(in interface{}) string {
	fullTypeName := fmt.Sprintf("%T", in)
	nameComponents := strings.Split(fullTypeName, ".")
	return nameComponents[len(nameComponents)-1]
}
