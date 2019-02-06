package sc

import (
	"strconv"
	"time"
)

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

func readTime(str *string, res *time.Time) error {
	var err error
	if len(*str) > 0 {
		var intTime int64
		intTime, err = strconv.ParseInt(*str, 10, 64)
		if err != nil {
			return err
		}
		*res = time.Unix(intTime, 0)
	}
	return nil
}
