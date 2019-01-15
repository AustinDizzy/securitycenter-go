package sc

import (
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

// Group is the application group located within SecurityCenter
type Group struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	UserCount    int       `json:"userCount"`
	Users        []User    `json:"users"`
	CreatedTime  time.Time `json:"createdTime"`
	ModifiedTime time.Time `json:"modifiedTime"`
}

func parseJSONToGroup(data *simplejson.Json) {
	for _, v := range []string{"id"} {
		key := strings.Split(v, ".")
		n, _ := strconv.Atoi(data.GetPath(key...).MustString())
		data.SetPath(key, n)
	}
	for _, v := range []string{"createdTime", "modifiedTime"} {
		n, _ := strconv.ParseInt(data.Get(v).MustString(), 10, 64)
		data.Set(v, time.Unix(n, 0))
	}
}
