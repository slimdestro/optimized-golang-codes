/**
	@ struct validation
	@ slimdestro
*/
package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string
	Age  int
}

func validateStruct(s interface{}) bool {
	// Get the type of the struct
	t := reflect.TypeOf(s)

	// Make sure it is a struct
	if t.Kind() != reflect.Struct {
		fmt.Println("Not a struct")
		return false
	}

	// Get the number of fields
	numFields := t.NumField()

	// Iterate over the fields
	for i := 0; i < numFields; i++ {
		// Get the field
		field := t.Field(i)

		// Get the field type
		fieldType := field.Type

		// Make sure the field type is supported
		if fieldType != reflect.TypeOf("") && fieldType != reflect.TypeOf(0) {
			fmt.Printf("Unsupported field type %s\n", fieldType)
			return false
		}
	}

	return true
}

func main() {
	p := Person{
		Name: "John Doe",
		Age:  42,
	}

	valid := validateStruct(p)
	fmt.Println(valid)
}