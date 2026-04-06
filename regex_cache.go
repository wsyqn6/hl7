package hl7

import (
	"regexp"
	"sync"
)

var regexCache = &sync.Map{}

func compileRegex(pattern string) *regexp.Regexp {
	if re, ok := regexCache.Load(pattern); ok {
		return re.(*regexp.Regexp)
	}

	re := regexp.MustCompile(pattern)
	regexCache.Store(pattern, re)
	return re
}

func ClearRegexCache() {
	regexCache = &sync.Map{}
}

func RegexCacheSize() int {
	count := 0
	regexCache.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}
