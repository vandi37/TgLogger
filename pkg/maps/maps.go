package maps

import (
	"sort"
	"strings"
)

func Find[V any](m map[string]V, s string, sep string) (V, bool) {
	if m == nil {
		return *new(V), false
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})

	for _, key := range keys {
		if strings.HasPrefix(s, key) && (len(key) == len(s) || (len(s) > len(key) && strings.HasPrefix(s[len(key):], sep))) {

			return m[key], true

		}
	}

	return *new(V), false
}
