package utilities

func CreateGenericMap[K int | float32 | string | []int | []float32 | []string](toConvert map[int]K) map[int]interface{} {
	m := make(map[int]interface{})

	for k, v := range toConvert {
		m[k] = v
	}
	return m
}
