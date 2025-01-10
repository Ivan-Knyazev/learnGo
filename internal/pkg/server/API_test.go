package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-storage/internal/pkg/storage"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScalarRouter(t *testing.T) {
	// Create strategy, storage, server
	strategy := storage.CreateJSONSaver("")
	storageObj, err := storage.NewStorage(strategy)
	if err != nil {
		log.Fatal("Error in creating new storage")
	}
	host := fmt.Sprintf("0.0.0.0:%d", 9000)
	s := NewServer(host, storageObj)

	// Testing router /scalar
	router := s.newAPI()
	w := httptest.NewRecorder()

	// Testing /scalar/set/test
	jsonData := `{"value": "4567", "ttl": 20}`
	body := bytes.NewBufferString(jsonData)
	req, _ := http.NewRequest(http.MethodPut, "/scalar/set/test", body)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "", w.Body.String())

	// Testing /scalar/get/test
	req, _ = http.NewRequest(http.MethodGet, "/scalar/get/test", nil)
	router.ServeHTTP(w, req)

	var response ScalarResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, storage.ScalarKindInt, response.ValueType)
	assert.Equal(t, "4567", response.Value)
}
