/* 
@ php substr_replace() in golang
@ slimdestro
*/
package main

import (
	"fmt"
	"strings"
)

func substr_replace(str string, replacement string, start int, length int) string {
	return strings.Replace(str[start:start+length], str[start:start+length], replacement, -1)
}

func main() {
	str := "Hello Destro"
	replacement := "Go"
	start := 6
	length := 5
	fmt.Println(substr_replace(str, replacement, start, length))
}