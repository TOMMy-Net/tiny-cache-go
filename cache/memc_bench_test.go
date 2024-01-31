package cache

import (
	"fmt"
	"testing"
	"time"
)

func insertXPreallocIntMap(x *Cache, b *testing.B) {

	for i := 0; i < 10000; i++ {
		x.Set(fmt.Sprint(i), i, 1*time.Hour)
		x.Get(fmt.Sprint(i))
	}
}

func BenchmarkCascadeCache(b *testing.B) {
	c := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		insertXPreallocIntMap(c, b)
	}
}
