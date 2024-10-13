package main

import (
	"fmt"
	"go-storage/internal/pkg/storage"
)

func main() {
	storageObj, err := storage.NewStorage()
	if err != nil {
		panic(err)
	}

	storageObj.Set("int", "1243232432")
	storageObj.Set("float", "-12432.2432")
	storageObj.Set("bool", "true")
	storageObj.Set("complex", "(2-3i)")
	storageObj.Set("string", "test_string-tatata rarara")

	// var data []any

	valueInt := storageObj.Get("int")
	valueFloat := storageObj.Get("float")
	valueBool := storageObj.Get("bool")
	valueComplex := storageObj.Get("complex")
	valueString := storageObj.Get("string")

	fmt.Println("\nTests of types:")
	fmt.Println(*valueInt, "type:", storageObj.GetKind("int"))
	fmt.Println(*valueFloat, "type:", storageObj.GetKind("float"))
	fmt.Println(*valueBool, "type:", storageObj.GetKind("bool"))
	fmt.Println(*valueComplex, "type:", storageObj.GetKind("complex"))
	fmt.Println(*valueString, "type:", storageObj.GetKind("string"))

	// fmt.Println("Tests:")
	// fmt.Println(valueInt, "type:", reflect.TypeOf(valueInt))
	// fmt.Println(valueFloat, "type:", reflect.TypeOf(valueFloat))
	// fmt.Println(valueBool, "type:", reflect.TypeOf(valueBool))
	// fmt.Println(valueComplex, "type:", reflect.TypeOf(valueComplex))
	// fmt.Println(valueString, "type:", reflect.TypeOf(valueString))
}
