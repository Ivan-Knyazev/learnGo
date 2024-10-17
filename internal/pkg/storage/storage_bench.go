package storage

// import "testing"

// type benchCase struct {
// 	name  string
// 	key   string
// 	value string
// 	kind  Kind
// }

// var bencCases = []benchCase{
// 	{"string_test_2", "hello", "world", "S"},
// 	{"string_test_1", "hello", "tetete", "S"},
// 	{"int_test", "int", "-12432", "D"},
// 	{"float_test", "float", "343232.24314323", "F"}, // ??? falue with value = 2321243232432.243234
// 	{"bool_test", "bool", "false", "B"},
// 	{"complex_test", "complex", "(2-3i)", "C"},
// }

// func BenchmarkGet(b *testing.B) {
// 	testStorage, err := NewStorage()
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, tCase := range bencCases {
// 		b.Run(tCase.name, func(bb *testing.B) {
// 			testStorage.Set(tCase.key, tCase.value)

// 			bb.ResetTimer()

// 			for i := 0; i < bb.N; i++ {
// 				testStorage.Get(tCase.key)
// 			}
// 		})
// 	}
// }

// type bench struct {
// 	name       string
// 	countCalls int
// }

// var cases = []bench{
// 	{"10", 10},
// 	{"100", 100},
// 	{"1000", 1000},
// 	{"5000", 5000},
// 	{"10000", 10000},
// 	{"100000", 100000},
// }
