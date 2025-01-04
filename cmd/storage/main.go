package main

import (
	"context"
	"fmt"
	"go-storage/internal/pkg/server"
	"go-storage/internal/pkg/storage"
	"go-storage/internal/pkg/utils"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
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

	if err = storage.ReadFromFile(storageObj, envs["path"]); err != nil {
		log.Println(err)
	}

	// Start cleen old data and save state in interval
	closeChan := make(chan struct{})
	interval, err := strconv.Atoi(envs["interval"])
	if err != nil {
		log.Println(err)
	}
	storageObj.StartScheduling(closeChan, int64(interval), storageObj, envs["path"])

	host := fmt.Sprintf("0.0.0.0:%s", envs["port"])
	s := server.NewServer(host, storageObj)

	// gin.SetMode(gin.ReleaseMode) // On release mode

	httpServer := s.StartServer()

	// Catching os exit signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //SIGINT (kill -2), SIGTERM(kill)
	// Waiting signal
	<-quit
	s.Storage.WriteLog("Shutdown Server ...")

	timeout, err := strconv.Atoi(envs["timeout"])
	if err != nil {
		log.Println(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(int64(timeout))*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// Stop scheduling
	close(closeChan)

	// Save state
	s.Storage.WriteLog("Save state of Storage ...")
	if err := storage.WriteToFile(storageObj, envs["path"]); err != nil {
		log.Fatal(err)
	}

	// Catching ctx.Done(). Timeout of <timeout> seconds
	<-ctx.Done()
	log.Printf("timeout of %d seconds", timeout)
	s.Storage.WriteLog("Server exiting")
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
