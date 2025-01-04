package storage

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Storage interface
type Storage interface {
	// Methods for scalar
	SetScalar(ttl int64, key string, val string) error
	GetScalar(key string) (string, bool, int64)
	GetScalarKind(key string) ScalarKind
	// Methods for dict
	SetDictFields(ttl int64, key string, elements ...string) (int, error) // HSET
	GetDictField(key string, field string) (ScalarValue, int64, error)    // HGET
	GetDict(key string) (map[string]ScalarValue, int64, error)
	// Methods for slice
	GetSlice(key string) ([]int, int64, error)
	LeftPushIntoSlice(ttl int64, key string, elements ...int) error        // LPUSH
	RightPushIntoSlice(ttl int64, key string, elements ...int) error       // RPUSH
	RightUniquePushIntoSlice(ttl int64, key string, elements ...int) error // RADDTOSET
	LeftPopFromSlice(key string, count ...int) (int, error)                // LPOP
	RightPopFromSlice(key string, count ...int) (int, error)               // RPOP
	SetSliceValue(key string, index int, element int) error                // LSET
	GetSliceValue(key string, index int) (int, int64, error)               // LGET
	// Methods for marshalling data for work with JSON
	LoadData(newData JsonStorage)
	ExportData() JsonStorage
	// Method for Scheduling
	StartScheduling(closeChan chan struct{}, shedulerInterval int64, storageObj Storage, JSONPath string)
	// Write logs
	WriteLog(info string)
	WriteLogWithParametr(info string, data any)
}

// type ScalarValue
type ScalarKind string

const (
	ScalarKindInt       ScalarKind = "D"
	ScalarKindString    ScalarKind = "S"
	ScalarKindUndefined ScalarKind = "UN"
)

type ScalarValue struct {
	ScalarValueType   ScalarKind `json:"ScalarValueType"`
	ScalarValueInt    int64      `json:"ScalarValueInt"`
	ScalarValueString string     `json:"ScalarValueString"`
}

// main type Value
type Kind string

const (
	KindScalar Kind = "SC"
	KindSlice  Kind = "SL"
	KindDict   Kind = "D"
)

type Value struct {
	ValueType Kind                   `json:"valueType"`
	Scalar    ScalarValue            `json:"scalar"`
	Slice     []int                  `json:"slice"`
	Dict      map[string]ScalarValue `json:"dict"`
	ExpiresAt int64                  `json:"expiresAt"`
}

// storage struct - implementation of Storage interface
type storage struct {
	data   map[string]Value
	Logger *zap.Logger
	mutex  sync.Mutex
}

// Create a new zap logger config
func newConfig() zap.Config {
	return zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout", "/tmp/go-storage-logs"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// Create a new storage
func NewStorage() (Storage, error) {

	config := newConfig()
	logger := zap.Must(config.Build())
	defer logger.Sync()

	// logger, err := zap.NewProduction()
	// if err != nil {
	// 	return Storage{}, err
	// }

	logger.Info("logger construction succeeded")
	logger.Info("created new storage")

	storage := &storage{
		data:   make(map[string]Value),
		Logger: logger,
		mutex:  sync.Mutex{},
	}
	return storage, nil
}

func (s *storage) StartScheduling(closeChan chan struct{}, shedulerInterval int64, storageObj Storage, JSONPath string) {
	interval := time.Duration(shedulerInterval) * time.Second
	go scheduler(s, closeChan, interval, storageObj, JSONPath)
	s.WriteLog("Start scheduling")
}

// For Marshalling
type JsonStorage struct {
	Data map[string]Value `json:"data"`
}

func (s *storage) LoadData(newData JsonStorage) {
	s.data = newData.Data
}

func (s *storage) ExportData() JsonStorage {
	return JsonStorage{
		Data: s.data,
	}
}

// For logging
func (s *storage) WriteLog(info string) {
	s.Logger.Info(fmt.Sprintf("[server] %s", info))
	defer s.Logger.Sync()
}

func (s *storage) WriteLogWithParametr(info string, data any) {
	s.Logger.Info(fmt.Sprintf("[server] %s", info), zap.Any("data", data))
	defer s.Logger.Sync()
}
