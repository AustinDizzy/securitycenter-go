package sc

func mergeData(data ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, set := range data {
		for k, v := range set {
			result[k] = v
		}
	}
	return result
}
