package helper

import (
	"testing"
)

func TestPrintLog2Hist(t *testing.T) {
	// 测试数据
	vals := []uint32{
		1,             // 2^0
		1 << 3,        // 2^3 = 8
		(1 << 7) + 10, // 2^7 + 10 = 138
		1 << 9,        // 2^9 = 512
		(1 << 10) + 5, // 2^10 + 5 = 1029
		1 << 4,        // 2^4 = 16
	}

	expected := `     qaq                 : count    distribution
         0 -> 1          : 1        |                                        |
         2 -> 3          : 8        |                                        |
         4 -> 7          : 138      |*****                                   |
         8 -> 15         : 512      |*******************                     |
        16 -> 31         : 1029     |****************************************|
        32 -> 63         : 16       |                                        |
`

	result := PrintLog2Hist(vals, "qaq")

	if result != expected {
		t.Errorf("PrintLog2Hist output mismatch\nExpected:\n%s\nGot:\n%s", expected, result)
	}
}

// 基准测试
func BenchmarkPrintLog2Hist(b *testing.B) {
	vals := make([]uint32, 64)
	for i := range vals {
		vals[i] = uint32(i * i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PrintLog2Hist(vals, "benchmark")
	}
}
