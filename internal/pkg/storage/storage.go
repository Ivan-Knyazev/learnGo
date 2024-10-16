package storage

import (
	"strconv"

	"go.uber.org/zap"
)

type value struct {
	valueType    uint8      // Type of field used (0 - valueInt, 1 - valueFloat, ...)
	valueInt     int64      // type 0
	valueFloat   float64    // type 1
	valueBool    bool       // type 2
	valueComplex complex128 // type 3
	valueString  string     // type 4
	// valueAny     any
}

type Storage struct {
	innerString map[string]value
	logger      *zap.Logger
}

func NewStorage() (Storage, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return Storage{}, err
	}

	defer logger.Sync()

	logger.Info("created new storage")

	return Storage{
		innerString: make(map[string]value),
		logger:      logger,
	}, nil
}

func (s Storage) Set(key string, val string) {

	s.logger.Info("key was set", zap.String("key", key), zap.Any("value", val))
	defer s.logger.Sync()

	// Check to int64
	valueInt, err := strconv.ParseInt(val, 10, 64)
	if err == nil {
		s.innerString[key] = value{valueType: 0, valueInt: valueInt}
		return
	}

	// Check to Float
	valueFloat, err := strconv.ParseFloat(val, 64)
	if err == nil {
		s.innerString[key] = value{valueType: 1, valueFloat: valueFloat}
		return
	}

	// Check to Bool
	valueBool, err := strconv.ParseBool(val)
	if err == nil {
		s.innerString[key] = value{valueType: 2, valueBool: valueBool}
		return
	}

	// Check to Complex
	valueComplex, err := strconv.ParseComplex(val, 128)
	if err == nil {
		s.innerString[key] = value{valueType: 3, valueComplex: valueComplex}
		return
	}

	// Is string
	s.innerString[key] = value{valueType: 4, valueString: val}
}

func (s Storage) Get(key string) *string {
	val, ok := s.get(key)
	if !ok {
		return nil
	}

	switch valueType := val.valueType; valueType {
	case 0:
		strInt := strconv.FormatInt(val.valueInt, 10)
		return &strInt
	case 1:
		strFloat := strconv.FormatFloat(val.valueFloat, 'f', -1, 64)
		return &strFloat
	case 2:
		strBool := strconv.FormatBool(val.valueBool)
		return &strBool
	case 3:
		strComplex := strconv.FormatComplex(val.valueComplex, 'f', -1, 64)
		return &strComplex
	// case 4:
	// 	return &val.valueString
	default:
		return &val.valueString
	}
}

func (s Storage) get(key string) (value, bool) {
	val, ok := s.innerString[key]
	if !ok {
		return value{}, false
	}

	return val, true
}

type Kind string

const (
	KindInt       Kind = "D" // type 0
	KindFloat     Kind = "F" // type 1
	KindBool      Kind = "B" // type 2
	KindComplex   Kind = "C" // type 3
	KindString    Kind = "S" // type 4
	KindUndefined Kind = "UN"
)

func (s Storage) GetKind(key string) Kind {
	value, ok := s.innerString[key]
	if !ok {
		return KindUndefined
	}

	switch valueType := value.valueType; valueType {
	case 0:
		return KindInt
	case 1:
		return KindFloat
	case 2:
		return KindBool
	case 3:
		return KindComplex
	case 4:
		return KindString
	default:
		return KindUndefined
	}
}
