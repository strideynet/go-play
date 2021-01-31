package v1

import (
	"testing"
)

func BenchmarkBM_Subscribe(b *testing.B) {
	bm := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bm.Subscribe()
	}
}
