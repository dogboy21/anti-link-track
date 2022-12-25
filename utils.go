package main

func mapKeys[K comparable, V any](m map[K]V) []K {
	keysResult := make([]K, 0, len(m))

	for key := range m {
		keysResult = append(keysResult, key)
	}

	return keysResult
}
