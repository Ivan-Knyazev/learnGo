package storage

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// --------
// For Dict
func (s *storage) SetDictFields(ttl int64, key string, elements ...string) (int, error) {
	_, ok := s.data[key]
	if ok && s.data[key].ValueType != KindDict {
		return 0, errors.New("requested field is not of type dict")
	}

	dict := make(map[string]ScalarValue)
	if ok {
		dict = s.data[key].Dict
	}

	keysCount := 0
	for index, element := range elements {
		if index%2 != 0 {
			keysCount++
			valueInt, err := strconv.ParseInt(element, 10, 64) // Check to int64
			if err != nil {                                    // Is string
				dict[elements[index-1]] = ScalarValue{
					ScalarValueType:   ScalarKindString,
					ScalarValueString: element,
				}
			} else { // Is int64
				dict[elements[index-1]] = ScalarValue{
					ScalarValueType: ScalarKindInt,
					ScalarValueInt:  valueInt,
				}
			}
		}
	}
	timeDuration := time.Duration(ttl) * time.Second
	s.data[key] = Value{
		ValueType: KindDict,
		Slice:     make([]int, 0),
		Dict:      dict,
		ExpiresAt: time.Now().Add(timeDuration).UnixMilli(),
	}
	s.Logger.Info(fmt.Sprintf("dict <%s> was set", key), zap.Any("keys-values", elements))
	defer s.Logger.Sync()
	return keysCount, nil
}

func (s *storage) GetDictField(key string, field string) (ScalarValue, int64, error) {
	dict, ok := s.data[key]
	if !ok {
		return ScalarValue{}, 0, errors.New("key not found")
	}
	if s.data[key].ValueType != KindDict {
		return ScalarValue{}, 0, errors.New("requested field is not of type dict")
	}

	value, ok := dict.Dict[field]
	if !ok {
		return ScalarValue{}, 0, fmt.Errorf("field %s not found in dict %s", field, key)
	}

	return value, s.data[key].ExpiresAt, nil
}

func (s *storage) GetDict(key string) (map[string]ScalarValue, int64, error) {
	value, ok := s.data[key]
	if !ok {
		return map[string]ScalarValue{}, 0, errors.New("key not found")
	}
	if s.data[key].ValueType != KindDict {
		return map[string]ScalarValue{}, 0, errors.New("requested field is not of type dict")
	}

	return value.Dict, s.data[key].ExpiresAt, nil
}
