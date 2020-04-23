package structgen

import "testing"

//BENCHMARK GetTokenAttributes
//Improvement history:
//------------------------------------
//	attempt	   |  time per loop
//------------------------------------
//  589063	      1962 ns/op
//  686304	      1792 ns/op (now)
//------------------------------------
func BenchmarkGetTokenAttributes(b *testing.B) {
	query := "id=1 && (division=engineering || division=finance)"
	for n := 0; n < b.N; n++ {
		getTokenAttributes(query)
	}
}
