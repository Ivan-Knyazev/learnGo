package storage

import (
	"errors"
	"slices"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Kind string

const (
	KindInt       Kind = "D"
	KindString    Kind = "S"
	KindUndefined Kind = "UN"
)

type Value struct {
	ValueType   Kind   `json:"valueType"`
	ValueInt    int64  `json:"valueInt"`
	ValueString string `json:"valueString"`
}

type Storage struct {
	innerValue map[string]Value
	innerSlice map[string][]int
	logger     *zap.Logger
}

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

	return Storage{
		innerValue: make(map[string]Value),
		innerSlice: make(map[string][]int),
		logger:     logger,
	}, nil
}

// For innerValue
func (s Storage) Set(key string, val string) {
	s.logger.Info("key was set", zap.String("key", key), zap.Any("value", val))
	defer s.logger.Sync()

	valueInt, err := strconv.ParseInt(val, 10, 64) // Check to int64
	if err != nil {                                // Is string
		s.innerValue[key] = Value{ValueType: KindString, ValueString: val}
	} else { // Is int64
		s.innerValue[key] = Value{ValueType: KindInt, ValueInt: valueInt}
	}
}

func (s Storage) Get(key string) (string, bool) {
	val, ok := s.get(key)
	if !ok {
		return "", false
	}

	switch valueType := val.ValueType; valueType {
	case KindInt:
		strInt := strconv.FormatInt(val.ValueInt, 10)
		return strInt, true
	case KindString:
		return val.ValueString, true
	default:
		return "", false
	}
}

func (s Storage) get(key string) (Value, bool) {
	val, ok := s.innerValue[key]
	return val, ok
}

func (s Storage) GetKind(key string) Kind {
	value, ok := s.innerValue[key]
	if !ok {
		return KindUndefined
	}
	return value.ValueType
}

// Func for testing innerSlice
func (s Storage) GetSlice(key string) []int {
	return s.innerSlice[key]
}

// For innerSlice
func (s Storage) LPUSH(key string, elements ...int) {
	slices.Reverse(elements)

	slice, ok := s.innerSlice[key]
	if !ok {
		s.innerSlice[key] = elements
		return
	}

	s.innerSlice[key] = slices.Concat(elements, slice)
}

func (s Storage) RPUSH(key string, elements ...int) {
	_, ok := s.innerSlice[key]
	if !ok {
		s.innerSlice[key] = elements
		return
	}

	s.innerSlice[key] = append(s.innerSlice[key], elements...)
}

func (s Storage) RADDTOSET(key string, elements ...int) {
	for _, element := range elements {
		if !slices.Contains(s.innerSlice[key], element) {
			s.innerSlice[key] = append(s.innerSlice[key], element)
		}
	}
}

func (s Storage) LPOP(key string, count ...int) int {
	_, ok := s.innerSlice[key]
	if !ok {
		return -1
	}

	if len(count) == 0 {
		return len(s.innerSlice[key])
	} else if len(count) == 1 {
		end := count[0]
		if end > 0 && end <= len(s.innerSlice[key]) {
			deleted := s.innerSlice[key][end-1]
			s.innerSlice[key] = slices.Delete(s.innerSlice[key], 0, end)
			return deleted
		} else if end > 0 && end > len(s.innerSlice[key]) {
			return len(s.innerSlice[key])
		} else {
			return -1
		}
	} else if len(count) == 2 {
		start := count[0]
		end := count[1]
		if start < 0 {
			start = len(s.innerSlice[key]) + start
		}
		if end < 0 {
			end = len(s.innerSlice[key]) + end
		}
		if end-start < 0 || start < 0 || end < 0 {
			return -1
		}

		if start >= 0 && start < len(s.innerSlice[key]) && end >= 0 && end < len(s.innerSlice[key]) {
			deleted := s.innerSlice[key][end]
			s.innerSlice[key] = slices.Delete(s.innerSlice[key], start, end+1)
			return deleted
		} else {
			return len(s.innerSlice[key]) - start
		}
	} else {
		return -1
	}
}

func (s Storage) RPOP(key string, count ...int) int {
	_, ok := s.innerSlice[key]
	if !ok {
		return -1
	}

	if len(count) == 0 {
		return len(s.innerSlice[key])
	} else if len(count) == 1 {
		offset := count[0]
		lenght := len(s.innerSlice[key])
		if offset > 0 && lenght-offset >= 0 {
			deleted := s.innerSlice[key][lenght-1]
			s.innerSlice[key] = slices.Delete(s.innerSlice[key], lenght-offset, lenght)
			return deleted
		} else if offset > 0 && lenght-offset < 0 {
			return len(s.innerSlice[key])
		} else {
			return -1
		}
	} else if len(count) == 2 {
		start := count[0]
		end := count[1]
		if start < 0 {
			start = len(s.innerSlice[key]) + start
		}
		if end < 0 {
			end = len(s.innerSlice[key]) + end
		}
		if end-start < 0 || start < 0 || end < 0 {
			return -1
		}
		// fmt.Println(start, end)

		if start >= 0 && start < len(s.innerSlice[key]) && end >= 0 && end < len(s.innerSlice[key]) {
			deleted := s.innerSlice[key][end]
			s.innerSlice[key] = slices.Delete(s.innerSlice[key], start, end+1)
			return deleted
		} else {
			return len(s.innerSlice[key]) - start
		}
	} else {
		return -1
	}
}

func (s Storage) LSET(key string, index int, element int) error {
	value, ok := s.innerSlice[key]
	if !ok {
		return errors.New("key not found")
	}
	if index < 0 || index >= len(value) {
		return errors.New("index out of range")
	}
	value[index] = element
	return nil
}

func (s Storage) LGET(key string, index int) (int, error) {
	value, ok := s.innerSlice[key]
	if !ok {
		return 0, errors.New("key not found")
	}
	if index < 0 || index >= len(value) {
		return 0, errors.New("index out of range")
	}
	return value[index], nil
}

// For Marshalling
type JsonStorage struct {
	InnerValue map[string]Value `json:"innerValue"`
	InnerSlice map[string][]int `json:"innerSlice"`
}

func (s *Storage) LoadData(data JsonStorage) {
	s.innerValue = data.InnerValue
	s.innerSlice = data.InnerSlice
}

func (s *Storage) ExportData() JsonStorage {
	return JsonStorage{
		InnerSlice: s.innerSlice,
		InnerValue: s.innerValue,
	}
}
