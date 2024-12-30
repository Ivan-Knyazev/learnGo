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

	envs := utils.GetEnvs(storageObj) // Comment to debug
	// For debug in VS Code
	// envs := make(map[string]string)
	// envs["path"] = "../../storage.json"
	// envs["port"] = "8090"

	if err = utils.ReadFromFile(storageObj, envs["path"]); err != nil {
		log.Println(err)
	}

	host := fmt.Sprintf("0.0.0.0:%s", envs["port"])
	s := server.NewServer(host, storageObj)
	if err := s.StartServer(); err != nil {
		log.Println(err)
	} else {
		defer s.Storage.WriteLog(fmt.Sprintf("server was started on %s", s.Host))
	}

	// storageObj.SetScalar("int", "1243232432")
	// storageObj.SetScalar("string", "test_string-tatata rarara")
	// storageObj.RPUSH("slice1", 1, 10, 3, 5, 8, 4, 10, 11)
	// storageObj.RPUSH("slice2")
	// storageObj.RPUSH("slice3", 1, 3, 5)
	// valueInt, ok := storageObj.GetScalar("int")
	// if !ok {
	// 	log.Println("invalid value at any key")
	// }
	// valueString, ok := storageObj.GetScalar("string")
	// if !ok {
	// 	log.Println("invalid value at any key")
	// }
	// nothingValue, ok := storageObj.GetScalar("nothing")
	// if !ok {
	// 	log.Println("invalid value at any key")
	// }
	// fmt.Println("\nTests of types:")
	// fmt.Println(valueInt, "type:", storageObj.GetScalarKind("int"))
	// fmt.Println(valueString, "type:", storageObj.GetScalarKind("string"))
	// fmt.Println(nothingValue, "type:", storageObj.GetScalarKind("nothing")+"\n")

	if err = utils.WriteToFile(storageObj, envs["path"]); err != nil {
		log.Fatal(err)
	}
}
