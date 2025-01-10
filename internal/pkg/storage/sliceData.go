package storage

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"go.uber.org/zap"
)

// ---------
// For Slice
func (s *storage) GetSlice(key string) ([]int, int64, error) {
	value, ok := s.data[key]
	if !ok {
		return []int{}, 0, errors.New("key not found")
	}
	if ok && s.data[key].ValueType != KindSlice {
		return []int{}, 0, errors.New("requested field is not of type slice")
	}

	return value.Slice, value.ExpiresAt, nil
}

// Inserts on the left into the list
func (s *storage) LeftPushIntoSlice(ttl int64, key string, elements ...int) error {
	slices.Reverse(elements)

	timeDuration := time.Duration(ttl) * time.Second

	value, ok := s.data[key]
	if !ok {
		s.data[key] = Value{
			ValueType: KindSlice,
			Slice:     append([]int{}, elements...),
			Dict:      make(map[string]ScalarValue),
			ExpiresAt: time.Now().Add(timeDuration).UnixMilli(),
		}
	} else {
		if s.data[key].ValueType != KindSlice {
			return errors.New("requested field is not of type slice")
		}
		s.data[key] = Value{
			ValueType: KindSlice,
			Slice:     slices.Concat(elements, value.Slice),
			Dict:      make(map[string]ScalarValue),
			ExpiresAt: time.Now().Add(timeDuration).UnixMilli(),
		}
	}
	s.Logger.Info(fmt.Sprintf("slice <%s> was set - left push", key), zap.Any("keys-values", elements))
	defer s.Logger.Sync()
	return nil
}

// Inserts on the right into the list
func (s *storage) RightPushIntoSlice(ttl int64, key string, elements ...int) error {
	timeDuration := time.Duration(ttl) * time.Second

	value, ok := s.data[key]
	if !ok {
		s.data[key] = Value{
			ValueType: KindSlice,
			Slice:     append([]int{}, elements...),
			Dict:      make(map[string]ScalarValue),
			ExpiresAt: time.Now().Add(timeDuration).UnixMilli(),
		}
	} else {
		if s.data[key].ValueType != KindSlice {
			return errors.New("requested field is not of type slice")
		}
		s.data[key] = Value{
			ValueType: KindSlice,
			Slice:     append(value.Slice, elements...),
			Dict:      make(map[string]ScalarValue),
			ExpiresAt: time.Now().Add(timeDuration).UnixMilli(),
		}
	}
	s.Logger.Info(fmt.Sprintf("slice <%s> was set - right push", key), zap.Any("keys-values", elements))
	defer s.Logger.Sync()
	return nil
}

// Inserts elements that are not yet in the list into the list on the right.
func (s *storage) RightUniquePushIntoSlice(ttl int64, key string, elements ...int) error {
	timeDuration := time.Duration(ttl) * time.Second
	ExpiresAt := time.Now().Add(timeDuration).UnixMilli()

	_, ok := s.data[key]
	if ok && s.data[key].ValueType != KindSlice {
		return errors.New("requested field is not of type slice")
	}

	for index, element := range elements {
		if !slices.Contains(s.data[key].Slice, element) {
			if !ok && index == 0 {
				s.data[key] = Value{
					ValueType: KindSlice,
					Slice:     append([]int{}, element),
					Dict:      make(map[string]ScalarValue),
					ExpiresAt: ExpiresAt,
				}
				continue
			}
			s.data[key] = Value{
				ValueType: KindSlice,
				Slice:     append(s.data[key].Slice, element),
				Dict:      make(map[string]ScalarValue),
				ExpiresAt: ExpiresAt,
			}
		}
	}
	s.Logger.Info(fmt.Sprintf("slice <%s> was set - right unique push", key), zap.Any("keys-values", elements))
	defer s.Logger.Sync()
	return nil
}

// Check slice key in storage
func (s *storage) checkSliceKey(key string) error {
	_, ok := s.data[key]
	if !ok {
		return errors.New("key not found")
	}
	if ok && s.data[key].ValueType != KindSlice {
		return errors.New("requested field is not of type slice")
	}

	return nil
}

// Deletes left element and return it
func (s *storage) LeftPopFromSlice(key string, count ...int) (int, error) {
	err := s.checkSliceKey(key)
	if err != nil {
		return -1, err
	}

	defer s.Logger.Sync()

	if len(count) == 0 {
		return len(s.data[key].Slice), nil
	} else if len(count) == 1 {
		end := count[0]
		if end > 0 && end <= len(s.data[key].Slice) {
			deleted := s.data[key].Slice[end-1]
			s.data[key] = Value{
				ValueType: KindSlice,
				Slice:     slices.Delete(s.data[key].Slice, 0, end),
				Dict:      make(map[string]ScalarValue),
				ExpiresAt: s.data[key].ExpiresAt,
			}
			s.Logger.Info(fmt.Sprintf("slice <%s> was set - left pop", key), zap.Any("last deleted element", deleted))
			return deleted, nil
		} else if end > 0 && end > len(s.data[key].Slice) {
			return len(s.data[key].Slice), nil
		} else {
			return -1, errors.New("error deleting an element")
		}
	} else if len(count) == 2 {
		start := count[0]
		end := count[1]
		if start < 0 {
			start = len(s.data[key].Slice) + start
		}
		if end < 0 {
			end = len(s.data[key].Slice) + end
		}
		if end-start < 0 || start < 0 || end < 0 {
			return -1, errors.New("error deleting an element")
		}

		if start >= 0 && start < len(s.data[key].Slice) && end >= 0 && end < len(s.data[key].Slice) {
			deleted := s.data[key].Slice[end]
			s.data[key] = Value{
				ValueType: KindSlice,
				Slice:     slices.Delete(s.data[key].Slice, start, end+1),
				Dict:      make(map[string]ScalarValue),
				ExpiresAt: s.data[key].ExpiresAt,
			}
			s.Logger.Info(fmt.Sprintf("slice <%s> was set - left pop", key), zap.Any("last deleted element", deleted))
			return deleted, nil
		} else {
			return len(s.data[key].Slice) - start, nil
		}
	} else {
		return -1, errors.New("error deleting an element")
	}
}

// Deletes right element and return it
func (s *storage) RightPopFromSlice(key string, count ...int) (int, error) {
	err := s.checkSliceKey(key)
	if err != nil {
		return -1, err
	}

	defer s.Logger.Sync()

	if len(count) == 0 {
		return len(s.data[key].Slice), nil
	} else if len(count) == 1 {
		offset := count[0]
		lenght := len(s.data[key].Slice)
		if offset > 0 && lenght-offset >= 0 {
			deleted := s.data[key].Slice[lenght-1]
			s.data[key] = Value{
				ValueType: KindSlice,
				Slice:     slices.Delete(s.data[key].Slice, lenght-offset, lenght),
				Dict:      make(map[string]ScalarValue),
				ExpiresAt: s.data[key].ExpiresAt,
			}
			s.Logger.Info(fmt.Sprintf("slice <%s> was set - right pop", key), zap.Any("last deleted element", deleted))
			return deleted, nil
		} else if offset > 0 && lenght-offset < 0 {
			return len(s.data[key].Slice), nil
		} else {
			return -1, errors.New("error deleting an element")
		}
	} else if len(count) == 2 {
		start := count[0]
		end := count[1]
		if start < 0 {
			start = len(s.data[key].Slice) + start
		}
		if end < 0 {
			end = len(s.data[key].Slice) + end
		}
		if end-start < 0 || start < 0 || end < 0 {
			return -1, errors.New("error deleting an element")
		}
		// fmt.Println(start, end)

		if start >= 0 && start < len(s.data[key].Slice) && end >= 0 && end < len(s.data[key].Slice) {
			deleted := s.data[key].Slice[end]
			s.data[key] = Value{
				ValueType: KindSlice,
				Slice:     slices.Delete(s.data[key].Slice, start, end+1),
				Dict:      make(map[string]ScalarValue),
				ExpiresAt: s.data[key].ExpiresAt,
			}
			s.Logger.Info(fmt.Sprintf("slice <%s> was set - right pop", key), zap.Any("last deleted element", deleted))
			return deleted, nil
		} else {
			return len(s.data[key].Slice) - start, nil
		}
	} else {
		return -1, errors.New("error deleting an element")
	}
}

// Sets value of element
func (s *storage) SetSliceValue(key string, index int, element int) error {
	value, ok := s.data[key]
	if !ok {
		return errors.New("key not found")
	}
	if s.data[key].ValueType != KindSlice {
		return errors.New("requested field is not of type slice")
	}
	if index < 0 || index >= len(value.Slice) {
		return errors.New("index out of range")
	}
	value.Slice[index] = element

	s.Logger.Info(fmt.Sprintf("slice <%s> was set - set %d element", key, index), zap.Any("value", element))
	defer s.Logger.Sync()
	return nil
}

// Gets value of element
func (s *storage) GetSliceValue(key string, index int) (int, int64, error) {
	value, ok := s.data[key]
	if !ok || s.data[key].ValueType != KindSlice {
		return 0, 0, errors.New("key not found")
	}
	if s.data[key].ValueType != KindSlice {
		return 0, 0, errors.New("requested field is not of type slice")
	}
	if index < 0 || index >= len(value.Slice) {
		return 0, 0, errors.New("index out of range")
	}
	return value.Slice[index], value.ExpiresAt, nil
}
