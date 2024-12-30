package storage

import (
	"errors"
	"fmt"
	"strconv"

	"go.uber.org/zap"
)

// ----------
// For Scalar
func (s *storage) SetScalar(key string, val string) error {
	_, ok := s.data[key]
	if ok && s.data[key].ValueType != KindScalar {
		return errors.New("requested field is not of type scalar")
	}

	valueInt, err := strconv.ParseInt(val, 10, 64) // Check to int64
	if err != nil {                                // Is string
		s.data[key] = Value{
			ValueType: KindScalar,
			Scalar:    ScalarValue{ScalarValueType: ScalarKindString, ScalarValueString: val},
			Slice:     make([]int, 0),
			Dict:      make(map[string]ScalarValue),
		}
	} else { // Is int64
		s.data[key] = Value{
			ValueType: KindScalar,
			Scalar:    ScalarValue{ScalarValueType: ScalarKindInt, ScalarValueInt: valueInt},
			Slice:     make([]int, 0),
			Dict:      make(map[string]ScalarValue),
		}
	}

	s.Logger.Info(fmt.Sprintf("scalar <%s> was set", key), zap.String("value", val))
	defer s.Logger.Sync()
	return nil
}

func (s *storage) GetScalar(key string) (string, bool) {
	val, ok := s.getValue(key)
	if !ok || val.ValueType != KindScalar {
		return "", false
	}

	switch valueType := val.Scalar.ScalarValueType; valueType {
	case ScalarKindInt:
		strInt := strconv.FormatInt(val.Scalar.ScalarValueInt, 10)
		return strInt, true
	case ScalarKindString:
		return val.Scalar.ScalarValueString, true
	default:
		return "", false
	}
}

func (s *storage) getValue(key string) (Value, bool) {
	val, ok := s.data[key]
	return val, ok
}

func (s *storage) GetScalarKind(key string) ScalarKind {
	value, ok := s.data[key]
	if !ok {
		return ScalarKindUndefined
	}
	return value.Scalar.ScalarValueType
}
