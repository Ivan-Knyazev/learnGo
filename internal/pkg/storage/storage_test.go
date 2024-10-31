package storage

import (
	"testing"
)

type testCaseValue struct {
	name  string
	key   string
	value string
	kind  Kind
}

var casesValue = []testCaseValue{
	{"string_test_1", "hello", "world", "S"},
	{"string_test_2", "hello", "tetete", "S"},
	{"int_test", "int", "-12432", "D"},
	{"float_test", "float", "343232.24314323", "S"}, // ??? falue with value = 2321243232432.243234
	{"bool_test", "bool", "false", "S"},
	{"complex_test", "complex", "(2-3i)", "S"},
}

func TestSetGetWithType(t *testing.T) {
	testStorage, err := NewStorage()
	if err != nil {
		panic(err)
	}

	for _, test := range casesValue {
		t.Run(test.name, func(t *testing.T) {
			testStorage.Set(test.key, test.value)

			value, ok := testStorage.Get(test.key)
			if !ok {
				t.Errorf("invalid value at key=%v", test.key)
			}

			if value != test.value {
				t.Errorf("two values should be the equal")
				t.Errorf("test.value=%s, got value=%s", test.value, value)
			}
			// assert.Equal(t, *value, test.value, "The two values should be the equal.")

			kind := testStorage.GetKind(test.key)

			if kind != test.kind {
				t.Errorf("two kinds should be the equal")
				t.Errorf("test.kind=%s, got kind=%s", test.kind, kind)
			}
			// assert.Equal(t, kind, test.kind, "The two kinds should be the equal.")
		})
	}
}

type testCasePUSH struct {
	name     string
	key      string
	elements []int
	check    []int
}

var casesSlicePUSH = []testCasePUSH{
	{"RPUSH", "test1", []int{1, 2, 3}, []int{1, 2, 3}},
	{"RPUSH", "test1", []int{4, 5, 6, 7}, []int{1, 2, 3, 4, 5, 6, 7}},
	{"LPUSH", "test2", []int{}, []int{}},
	{"LPUSH", "test2", []int{1, 2, 3}, []int{3, 2, 1}},
	{"LPUSH", "test2", []int{4, 5, 6, 7}, []int{7, 6, 5, 4, 3, 2, 1}},
	{"RADDTOSET", "test3", []int{}, []int{}},
	{"RADDTOSET", "test3", []int{1, 2, 3}, []int{1, 2, 3}},
	{"RADDTOSET", "test3", []int{3, 5, 8, 4, 8, 10, 11}, []int{1, 2, 3, 5, 8, 4, 10, 11}},
}

func TestLPUSH(t *testing.T) {
	testStorage, err := NewStorage()
	if err != nil {
		panic(err)
	}

	for _, test := range casesSlicePUSH {
		t.Run(test.name, func(t *testing.T) {
			if test.name == "LPUSH" {
				testStorage.LPUSH(test.key, test.elements...)
			} else if test.name == "RPUSH" {
				testStorage.RPUSH(test.key, test.elements...)
			} else if test.name == "RADDTOSET" {
				testStorage.RADDTOSET(test.key, test.elements...)
			}

			slice := testStorage.GetSlice(test.key)

			if len(slice) < len(test.check) {
				t.Errorf("len of slice is less than required")
			} else if len(slice) > len(test.check) {
				t.Errorf("len of slice is more than required")
			}

			for index, value := range slice {
				if value != test.check[index] {
					t.Errorf("incorrect: value=%d at index=%d", value, index)
					t.Errorf("correct: value=%d at index=%d (from test)", test.check[index], index)
					t.Errorf("correct slice is %v really slice is %v", test.check, slice)
					break
				}
			}
			// assert.Equal(t, *value, test.value, "The two values should be the equal.")
		})
	}
}

type testCasePOP struct {
	name        string
	key         string
	arg1        any
	arg2        any
	checkAnswer int
	checkSlice  []int
}

var casesSlicePOP = []testCasePOP{
	{"LPOP", "test1", nil, nil, 8, []int{1, 2, 3, 5, 8, 4, 10, 11}},
	{"LPOP", "test1", 2, nil, 2, []int{3, 5, 8, 4, 10, 11}},
	{"LPOP", "test1", 10, nil, 6, []int{3, 5, 8, 4, 10, 11}},
	{"LPOP", "test1", -10, nil, -1, []int{3, 5, 8, 4, 10, 11}},
	{"LPOP", "test1", 3, 2, -1, []int{3, 5, 8, 4, 10, 11}},
	{"LPOP", "test1", 2, 3, 4, []int{3, 5, 10, 11}},
	{"LPOP", "test1", 2, 2, 10, []int{3, 5, 11}},
	{"LPOP", "test1", 1, 20, 2, []int{3, 5, 11}},
	{"LPOP", "test1", 1, -2, 5, []int{3, 11}},
	{"LPOP", "test1", -2, -1, 11, []int{}},
	{"LPOP", "test3", 1, -2, -1, []int{}},

	{"RPOP", "test2", nil, nil, 8, []int{1, 2, 3, 5, 8, 4, 10, 11}},
	{"RPOP", "test2", 2, nil, 11, []int{1, 2, 3, 5, 8, 4}},
	{"RPOP", "test2", 10, nil, 6, []int{1, 2, 3, 5, 8, 4}},
	{"RPOP", "test2", -10, nil, -1, []int{1, 2, 3, 5, 8, 4}},
	{"RPOP", "test2", 3, 2, -1, []int{1, 2, 3, 5, 8, 4}},
	{"RPOP", "test2", 2, 3, 5, []int{1, 2, 8, 4}},
	{"RPOP", "test2", 2, 2, 8, []int{1, 2, 4}},
	{"RPOP", "test2", 1, 20, 2, []int{1, 2, 4}},
	{"RPOP", "test2", 1, -2, 2, []int{1, 4}},
	{"RPOP", "test2", -2, -1, 4, []int{}},
	{"RPOP", "test3", 1, -2, -1, []int{}},
}

const (
	LPOP = iota
	RPOP
)

func checkPOP(testStorage Storage, test testCasePOP, t *testing.T, popType int) {
	var result int
	if test.arg1 == nil && test.arg2 == nil {
		if popType == RPOP {
			result = testStorage.RPOP(test.key)
		} else if popType == LPOP {
			result = testStorage.LPOP(test.key)
		}
		if test.checkAnswer != result {
			t.Errorf("incorrect return value=%v", result)
		}
	} else if test.arg1 != nil && test.arg2 == nil {
		if popType == RPOP {
			result = testStorage.RPOP(test.key, test.arg1.(int))
		} else if popType == LPOP {
			result = testStorage.LPOP(test.key, test.arg1.(int))
		}
	} else {
		if popType == RPOP {
			result = testStorage.RPOP(test.key, test.arg1.(int), test.arg2.(int))
		} else if popType == LPOP {
			result = testStorage.LPOP(test.key, test.arg1.(int), test.arg2.(int))
		}
	}
	if result != test.checkAnswer {
		t.Errorf("incorrect return value=%v", result)
	}
	slice := testStorage.GetSlice(test.key)
	if len(slice) != len(test.checkSlice) {
		t.Errorf("len of slices is not equal")
	}
	for index, value := range slice {
		if value != test.checkSlice[index] {
			t.Errorf("incorrect: value=%d at index=%d", value, index)
			t.Errorf("correct: value=%d at index=%d (from test)", test.checkSlice[index], index)
			t.Errorf("correct slice is %v really slice is %v", test.checkSlice, slice)
			break
		}
	}
}
func TestPOP(t *testing.T) {
	testStorage, err := NewStorage()
	if err != nil {
		panic(err)
	}
	testStorage.RPUSH("test1", 1, 2, 3, 5, 8, 4, 10, 11)
	testStorage.RPUSH("test2", 1, 2, 3, 5, 8, 4, 10, 11)

	for _, test := range casesSlicePOP {
		t.Run(test.name, func(t *testing.T) {
			if test.name == "LPOP" {
				checkPOP(testStorage, test, t, LPOP)
			} else if test.name == "RPOP" {
				checkPOP(testStorage, test, t, RPOP)
			}
		})
	}
}

type testCaseLSET struct {
	name       string
	key        string
	arg1       int
	arg2       int
	checkSlice []int
}

var casesSliceLSET = []testCaseLSET{
	{"LSET", "test1", 1, 10, []int{1, 10, 3, 5, 8, 4, 10, 11}},
	{"LSET", "test1", 0, 0, []int{0, 10, 3, 5, 8, 4, 10, 11}},
	{"LSET", "test1", 7, 8, []int{0, 10, 3, 5, 8, 4, 10, 8}},
	{"LSET", "test1", 10, 10, []int{0, 10, 3, 5, 8, 4, 10, 8}},
	{"LSET", "test1", -1, 10, []int{0, 10, 3, 5, 8, 4, 10, 8}},
	{"LSET", "test2", 1, 10, []int{}},
}

func TestLSET(t *testing.T) {
	testStorage, err := NewStorage()
	if err != nil {
		panic(err)
	}
	testStorage.RPUSH("test1", 1, 2, 3, 5, 8, 4, 10, 11)

	for _, test := range casesSliceLSET {
		t.Run(test.name, func(t *testing.T) {
			err := testStorage.LSET(test.key, test.arg1, test.arg2)
			if err != nil {
				t.Logf("error=%v", err)
			}
			slice := testStorage.GetSlice(test.key)
			if len(slice) != len(test.checkSlice) {
				t.Errorf("len of slices is not equal")
			}
			for index, value := range slice {
				if value != test.checkSlice[index] {
					t.Errorf("incorrect: value=%d at index=%d", value, index)
					t.Errorf("correct: value=%d at index=%d (from test)", test.checkSlice[index], index)
					t.Errorf("correct slice is %v really slice is %v", test.checkSlice, slice)
					break
				}
			}
		})
	}
}

type testCaseLGET struct {
	name       string
	key        string
	arg        int
	checkValue int
	checkSlice []int
}

var casesSliceLGET = []testCaseLGET{
	{"LGET", "test1", 1, 2, []int{1, 2, 3, 5, 8, 4, 10, 11}},
	{"LGET", "test1", 0, 1, []int{1, 2, 3, 5, 8, 4, 10, 11}},
	{"LGET", "test1", 7, 11, []int{1, 2, 3, 5, 8, 4, 10, 11}},
	{"LGET", "test1", 10, 0, []int{1, 2, 3, 5, 8, 4, 10, 11}},
	{"LGET", "test1", -1, 0, []int{1, 2, 3, 5, 8, 4, 10, 11}},
	{"LGET", "test2", 1, 0, []int{}},
}

func TestLLGET(t *testing.T) {
	testStorage, err := NewStorage()
	if err != nil {
		panic(err)
	}
	testStorage.RPUSH("test1", 1, 2, 3, 5, 8, 4, 10, 11)

	for _, test := range casesSliceLGET {
		t.Run(test.name, func(t *testing.T) {
			value, err := testStorage.LGET(test.key, test.arg)
			if err != nil {
				t.Logf("error=%v", err)
			}
			if value != test.checkValue {
				t.Errorf("incorrect return value=%v", value)
			}
			slice := testStorage.GetSlice(test.key)
			if len(slice) != len(test.checkSlice) {
				t.Errorf("len of slices is not equal")
			}
			for index, value := range slice {
				if value != test.checkSlice[index] {
					t.Errorf("incorrect: value=%d at index=%d", value, index)
					t.Errorf("correct: value=%d at index=%d (from test)", test.checkSlice[index], index)
					t.Errorf("correct slice is %v really slice is %v", test.checkSlice, slice)
					break
				}
			}
		})
	}
}

func BenchmarkGet(b *testing.B) {
	testStorage, err := NewStorage()
	if err != nil {
		panic(err)
	}

	for _, tCase := range casesValue {
		b.Run(tCase.name, func(bb *testing.B) {

			testStorage.Set(tCase.key, tCase.value)

			var value string
			bb.ResetTimer()

			for i := 0; i < bb.N; i++ {
				value, _ = testStorage.Get(tCase.key)
			}

			if value == tCase.value {
			}
		})
	}
}

func BenchmarkSet(b *testing.B) {
	testStorage, err := NewStorage()
	if err != nil {
		panic(err)
	}

	for _, tCase := range casesValue {
		b.Run(tCase.name, func(bb *testing.B) {

			// testStorage.logger.Info(tCase.name)
			bb.ResetTimer()

			for i := 0; i < bb.N; i++ {
				testStorage.Set(tCase.key, tCase.value)
			}
		})
	}
}

func BenchmarkSetGet(b *testing.B) {
	testStorage, err := NewStorage()
	if err != nil {
		panic(err)
	}

	for _, tCase := range casesValue {
		b.Run(tCase.name, func(bb *testing.B) {

			bb.ResetTimer()

			for i := 0; i < bb.N; i++ {
				testStorage.Set(tCase.key, tCase.value)
				testStorage.Get(tCase.key)
			}
		})
	}
}
