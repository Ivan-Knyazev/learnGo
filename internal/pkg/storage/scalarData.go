package storage

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// ----------
// For Scalar
func (s *storage) SetScalar(ttl int64, key string, val string) error {
	_, ok := s.data[key]
	if ok && s.data[key].ValueType != KindScalar {
		return errors.New("requested field is not of type scalar")
	}

	timeDuration := time.Duration(ttl) * time.Second

	valueInt, err := strconv.ParseInt(val, 10, 64) // Check to int64
	if err != nil {                                // Is string
		s.data[key] = Value{
			ValueType: KindScalar,
			Scalar:    ScalarValue{ScalarValueType: ScalarKindString, ScalarValueString: val},
			Slice:     make([]int, 0),
			Dict:      make(map[string]ScalarValue),
			ExpiresAt: time.Now().Add(timeDuration).UnixMilli(),
		}
	} else { // Is int64
		s.data[key] = Value{
			ValueType: KindScalar,
			Scalar:    ScalarValue{ScalarValueType: ScalarKindInt, ScalarValueInt: valueInt},
			Slice:     make([]int, 0),
			Dict:      make(map[string]ScalarValue),
			ExpiresAt: time.Now().Add(timeDuration).UnixMilli(),
		}
	}

	s.Logger.Info(fmt.Sprintf("scalar <%s> was set", key), zap.String("value", val))
	defer s.Logger.Sync()
	return nil
}

func (s *storage) GetScalar(key string) (string, bool, int64) {
	val, ok := s.getValue(key)
	if !ok || val.ValueType != KindScalar {
		return "", false, val.ExpiresAt
	}

	switch valueType := val.Scalar.ScalarValueType; valueType {
	case ScalarKindInt:
		strInt := strconv.FormatInt(val.Scalar.ScalarValueInt, 10)
		return strInt, true, val.ExpiresAt
	case ScalarKindString:
		return val.Scalar.ScalarValueString, true, val.ExpiresAt
	default:
		return "", false, val.ExpiresAt
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
