package hashs

import "testing"

const TestEncodingString = "benchmark_test_string"

func BenchmarkMD5(b *testing.B) {
    for i := 0; i < b.N; i++ {
        MD5(TestEncodingString)
    }
}

func BenchmarkMururHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
        MurmurHash(TestEncodingString)
    }
}

/*
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkMD5
BenchmarkMD5-12          	 4999347	       218.8 ns/op
BenchmarkMururHash
BenchmarkMururHash-12    	 6231925	       181.5 ns/op
*/
