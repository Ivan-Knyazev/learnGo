package main

import (
	"fmt"
	"go-storage/internal/pkg/server"
	"go-storage/internal/pkg/storage"
	"go-storage/internal/pkg/utils"
	"log"
)

func main() {
	storageObj, err := storage.NewStorage()
	if err != nil {
		log.Fatal(err)
	}

	envs := utils.GetEnvs(&storageObj)

	if err = utils.ReadFromFile(&storageObj, envs["path"]); err != nil {
		log.Println(err)
	}

	host := fmt.Sprintf("0.0.0.0:%s", envs["port"])
	s := server.NewServer(host, &storageObj)
	s.StartServer()

	// storageObj.Set("int", "1243232432")
	// storageObj.Set("string", "test_string-tatata rarara")
	// storageObj.RPUSH("slice1", 1, 10, 3, 5, 8, 4, 10, 11)
	// storageObj.RPUSH("slice2")
	// storageObj.RPUSH("slice3", 1, 3, 5)
	// valueInt, ok := storageObj.Get("int")
	// if !ok {
	// 	log.Println("invalid value at any key")
	// }
	// valueString, ok := storageObj.Get("string")
	// if !ok {
	// 	log.Println("invalid value at any key")
	// }
	// nothingValue, ok := storageObj.Get("nothing")
	// if !ok {
	// 	log.Println("invalid value at any key")
	// }
	// fmt.Println("\nTests of types:")
	// fmt.Println(valueInt, "type:", storageObj.GetKind("int"))
	// fmt.Println(valueString, "type:", storageObj.GetKind("string"))
	// fmt.Println(nothingValue, "type:", storageObj.GetKind("nothing"))

	// * Tests of innerSlice was added in storage_test.go

	if err = utils.WriteToFile(&storageObj, envs["path"]); err != nil {
		log.Fatal(err)
	}
}
