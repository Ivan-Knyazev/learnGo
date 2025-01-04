package storage

import (
	"fmt"
	"log"
	"time"
)

func scheduler(storage *storage, closeChan chan struct{}, interval time.Duration, storageObj Storage, JSONPath string) {
	for {
		select {
		case <-closeChan:
			return
		case <-time.After(interval):
			clean(storage)
			saveState(storageObj, JSONPath)
		}
	}
}

func clean(storage *storage) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	for key, value := range storage.data {
		if time.Now().UnixMilli() >= value.ExpiresAt {
			delete(storage.data, key)
			storage.WriteLog(fmt.Sprintf("Delete Value with key=%s", key))
		}
	}
}

func saveState(storageObj Storage, path string) {
	if err := WriteToFile(storageObj, path); err != nil {
		log.Fatal(err)
	}
}
