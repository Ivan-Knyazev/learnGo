package storage

import (
	"fmt"
	"time"
)

func scheduler(storage *storage, closeChan chan struct{}, interval time.Duration, storageObj Storage) {
	for {
		select {
		case <-closeChan:
			return
		case <-time.After(interval):
			clean(storage)
			storageObj.SaveData()
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
