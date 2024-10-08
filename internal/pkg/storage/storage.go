package storage

import (
	"strconv"
)

type Value struct {
	valueType    uint8      // Type of field used (0 - valueInt, 1 - valueFloat, ...)
	valueInt     int64      // type 0
	valueFloat   float64    // type 1
	valueBool    bool       // type 2
	valueComplex complex128 // type 3
	valueString  string     // type 4
	// valueAny     any
}

type Storage struct {
	innerString map[string]Value
	// logger      *zap.Logger
}

func NewStorage() (*Storage, error) {
	// logger, err := zap.NewProduction()
	// if err != nil {
	// 	return &Storage{}, err
	// }

	// defer logger.Sync()
	// logger.Info("created new storage")

	return &Storage{
		innerString: make(map[string]Value),
	}, nil
}

func (s *Storage) Set(key string, value string) {
	// s.logger.Info("key was set")
	// s.logger.Sync()

	// Check to int64
	valueInt, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		s.innerString[key] = Value{valueType: 0, valueInt: valueInt}
		return
	}

	// Check to Float
	valueFloat, err := strconv.ParseFloat(value, 64)
	if err == nil {
		s.innerString[key] = Value{valueType: 1, valueFloat: valueFloat}
		return
	}

	// Check to Bool
	valueBool, err := strconv.ParseBool(value)
	if err == nil {
		s.innerString[key] = Value{valueType: 2, valueBool: valueBool}
		return
	}

	// Check to Complex
	valueComplex, err := strconv.ParseComplex(value, 128)
	if err == nil {
		s.innerString[key] = Value{valueType: 3, valueComplex: valueComplex}
		return
	}

	// Is string
	s.innerString[key] = Value{valueType: 4, valueString: value}
}

func (s *Storage) Get(key string) any {
	value, ok := s.innerString[key]
	if !ok {
		return nil
	}

	switch valueType := value.valueType; valueType {
	case 0:
		return value.valueInt
	case 1:
		return value.valueFloat
	case 2:
		return value.valueBool
	case 3:
		return value.valueComplex
	default:
		return value.valueString
	}
}
