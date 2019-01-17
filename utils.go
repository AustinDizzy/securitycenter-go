package sc

import "strconv"

func mergeData(data ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, set := range data {
		for k, v := range set {
			result[k] = v
		}
	}
	return result
}

func readBool(str *string, res *bool) error {
	var err error
	if len(*str) > 0 {
		*res, err = strconv.ParseBool(*str)
	}
	return err
}
