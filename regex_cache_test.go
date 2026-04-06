package hl7

import (
	"regexp"
	"sync"
	"testing"
)

func TestCompileRegex(t *testing.T) {
	pattern := `^\d{8}$`

	re1 := compileRegex(pattern)
	re2 := compileRegex(pattern)

	if re1 != re2 {
		t.Error("compileRegex should return the same instance for the same pattern")
	}

	if !re1.MatchString("20240115") {
		t.Error("regex should match date format")
	}

	if re1.MatchString("invalid") {
		t.Error("regex should not match invalid format")
	}
}

func TestRegexCacheConcurrent(t *testing.T) {
	pattern := `^[A-Z]{2,3}$`

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			re := compileRegex(pattern)
			if !re.MatchString("ABC") {
				t.Error("regex should match")
			}
		}()
	}
	wg.Wait()
}

func TestClearRegexCache(t *testing.T) {
	initialSize := RegexCacheSize()

	pattern := `^test$`
	compileRegex(pattern)

	expectedSize := initialSize + 1
	if RegexCacheSize() != expectedSize {
		t.Errorf("expected cache size %d, got %d", expectedSize, RegexCacheSize())
	}

	ClearRegexCache()

	if RegexCacheSize() != 0 {
		t.Errorf("expected cache size 0 after clear, got %d", RegexCacheSize())
	}
}

func TestRegexCacheSize(t *testing.T) {
	ClearRegexCache()

	if RegexCacheSize() != 0 {
		t.Errorf("expected empty cache, got %d", RegexCacheSize())
	}

	compileRegex(`^pattern1$`)
	compileRegex(`^pattern2$`)
	compileRegex(`^pattern1$`) // duplicate

	if RegexCacheSize() != 2 {
		t.Errorf("expected cache size 2, got %d", RegexCacheSize())
	}
}

func BenchmarkCompileRegexCached(b *testing.B) {
	pattern := `^\d{8}$`
	compileRegex(pattern) // warm up cache

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		compileRegex(pattern)
	}
}

func BenchmarkCompileRegexUncached(b *testing.B) {
	pattern := `^\d{8}$`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		regexp.MustCompile(pattern)
	}
}
