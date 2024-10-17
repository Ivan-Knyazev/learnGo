package storage

import (
	"testing"
)

type testCase struct {
	name  string
	key   string
	value string
	kind  Kind
}

var cases = []testCase{
	{"string_test_1", "hello", "world", "S"},
	{"string_test_2", "hello", "tetete", "S"},
	{"int_test", "int", "-12432", "D"},
	{"float_test", "float", "343232.24314323", "F"}, // ??? falue with value = 2321243232432.243234
	{"bool_test", "bool", "false", "B"},
	{"complex_test", "complex", "(2-3i)", "C"},
}

func TestSetGetWithType(t *testing.T) {
	testStorage, err := NewStorage()
	if err != nil {
		panic(err)
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			testStorage.Set(test.key, test.value)

			value := testStorage.Get(test.key)

			if *value != test.value {
				t.Errorf("two values should be the equal")
				t.Errorf("test.value=%s, got value=%s", test.value, *value)
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

func BenchmarkGet(b *testing.B) {
	testStorage, err := NewStorage()
	if err != nil {
		panic(err)
	}

	for _, tCase := range cases {
		b.Run(tCase.name, func(bb *testing.B) {

			testStorage.Set(tCase.key, tCase.value)

			bb.ResetTimer()

			for i := 0; i < bb.N; i++ {
				testStorage.Get(tCase.key)
			}
		})
	}
}

func BenchmarkSet(b *testing.B) {
	testStorage, err := NewStorage()
	if err != nil {
		panic(err)
	}

	for _, tCase := range cases {
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

	for _, tCase := range cases {
		b.Run(tCase.name, func(bb *testing.B) {

			bb.ResetTimer()

			for i := 0; i < bb.N; i++ {
				testStorage.Set(tCase.key, tCase.value)
				testStorage.Get(tCase.key)
			}
		})
	}
}
