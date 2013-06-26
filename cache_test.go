package cache

import (
	"runtime"
	"sync"
	"testing"
)

func init() {
	// This should probably be varied to get to the proper scaling behavior
	runtime.GOMAXPROCS(12)
}

func TestSimpleMappingToIndex(t *testing.T) {
	str := "abcdefghijklmnopqrstuvxyzABCDEFGHIJKLMNOPQRSTUVXYZ1234567890"

	for _, chr := range str {
		h := indexFor(hash(string(chr), 0))
		if h < 0 || h > 15 {
			t.Fatalf("Out of bounds index detected, got '%v' for '%v'\n", h, chr)
		}
	}
}

func TestSimpleHashing(t *testing.T) {

	longString := "A" + "PAN"
	shortString := "A"

	h1 := hash(longString, 0)
	h2 := hash(shortString, 0)

	if h1 != h2 {
		t.Fatalf("Unequal hashing detected! got '%v' for '%v' and '%v' for '%v'", h1, longString, h2, shortString)
	}
}

func TestSimpleSetAndGetStandardCache(t *testing.T) {
	key := "apan"
	val := "bapan"

	c := NewCache()
	c.Set(key, val)

	v := c.Get(key)

	if v != val {
		t.Fatalf("Invalid value '%v' for key '%v', expected '%v'\n", v, key, val)
	}
}

func TestSimpleSetAndGetConcurrentCache(t *testing.T) {
	key := "apan"
	val := "bapan"

	c := NewConcurrentCache()
	c.Set(key, val)

	v := c.Get(key)

	if v != val {
		t.Fatalf("Invalid value '%v' for key '%v', expected '%v'\n", v, key, val)
	}
}

func BenchmarkHash(b *testing.B) {
	str := "APAN"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash(str, 0)
	}
}

func BenchmarkHashInt(b *testing.B) {
	v := 1
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hashInt(v)
	}
}

func BenchmarkIndexCalculation(b *testing.B) {
	v := 1
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h := hashInt(v)
		indexFor(h)
	}
}

func BenchmarkConcurrentCacheSet(b *testing.B) {
	str := "abcdefghijklmnopqrstuvxyzABCDEFGHIJKLMNOPQRSTUVXYZ1234567890"
	c := NewConcurrentCache()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, chr := range str {
			key := string(chr)
			c.Set(key+key, str)
		}
	}
}

func BenchmarkStandardCacheSet(b *testing.B) {
	str := "abcdefghijklmnopqrstuvxyzABCDEFGHIJKLMNOPQRSTUVXYZ1234567890"
	c := NewCache()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, chr := range str {
			key := string(chr)
			c.Set(key+key, str)
		}
	}
}

func BenchmarkStandardCacheConcurrentSet10(b *testing.B) {
	c := NewCache()
	b.ResetTimer()
	runCacheConcurrently(10, c, b)
}

func BenchmarkStandardCacheConcurrentSet20(b *testing.B) {
	c := NewCache()
	b.ResetTimer()
	runCacheConcurrently(20, c, b)
}

func BenchmarkStandardCacheConcurrentSet40(b *testing.B) {
	c := NewCache()
	b.ResetTimer()
	runCacheConcurrently(40, c, b)
}

func BenchmarkStandardCacheConcurrentSet80(b *testing.B) {
	c := NewCache()
	b.ResetTimer()
	runCacheConcurrently(80, c, b)
}

func BenchmarkStandardCacheConcurrentSet160(b *testing.B) {
	c := NewCache()
	b.ResetTimer()
	runCacheConcurrently(160, c, b)
}

func BenchmarkConcurrentCacheConcurrentSet10(b *testing.B) {
	c := NewConcurrentCache()
	b.ResetTimer()
	runCacheConcurrently(10, c, b)
}

func BenchmarkConcurrentCacheConcurrentSet20(b *testing.B) {
	c := NewConcurrentCache()
	b.ResetTimer()
	runCacheConcurrently(20, c, b)
}

func BenchmarkConcurrentCacheConcurrentSet40(b *testing.B) {
	c := NewConcurrentCache()
	b.ResetTimer()
	runCacheConcurrently(40, c, b)
}

func BenchmarkConcurrentCacheConcurrentSet80(b *testing.B) {
	c := NewConcurrentCache()
	b.ResetTimer()
	runCacheConcurrently(80, c, b)
}

func BenchmarkConcurrentCacheConcurrentSet160(b *testing.B) {
	c := NewConcurrentCache()
	b.ResetTimer()
	runCacheConcurrently(160, c, b)
}

func runCacheConcurrently(cnt int, c Cache, b *testing.B) {
	str := "abcdefghijklmnopqrstuvxyzABCDEFGHIJKLMNOPQRSTUVXYZ1234567890"

	for i := 0; i < b.N; i++ {
		var wgStart sync.WaitGroup
		var wg sync.WaitGroup

		wgStart.Add(1)
		wg.Add(cnt)

		for i := 0; i < cnt; i++ {
			go func() {
				wgStart.Wait()
				for _, chr := range str {
					key := string(chr)
					c.Set(key+key, str)
				}
				wg.Done()
			}()
		}
		wgStart.Done()
		wg.Wait()
	}
}

func BenchmarkSimpleMapSet(b *testing.B) {
	str := "abcdefghijklmnopqrstuvxyzABCDEFGHIJKLMNOPQRSTUVXYZ1234567890"
	c := make(map[string]string)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, chr := range str {
			key := string(chr)
			c[key] = str
		}
	}
}
