package main

import (
	"fmt"
	"go-storage/internal/pkg/storage"
	"log"
	"reflect"
)

func main() {
	storageObj, err := storage.NewStorage()
	if err != nil {
		log.Fatal(err)
		return
	}

	(*storageObj).Set("int", "1243232432")
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

	fmt.Println("Tests:")
	fmt.Println(valueInt, "type:", reflect.TypeOf(valueInt))
	fmt.Println(valueFloat, "type:", reflect.TypeOf(valueFloat))
	fmt.Println(valueBool, "type:", reflect.TypeOf(valueBool))
	fmt.Println(valueComplex, "type:", reflect.TypeOf(valueComplex))
	fmt.Println(valueString, "type:", reflect.TypeOf(valueString))
}
