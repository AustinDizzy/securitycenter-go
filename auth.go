package sc

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	sessionRgx = regexp.MustCompile(`(?:TNS_SESSIONID=)(.*?)(?:;)`)
)

// HasAuth returns true if the SC multiplexor has valid authentication information
func (sc *SC) HasAuth() bool {
	if session, ok := sc.Keys["session"]; ok && len(session) > 0 {
		if tokenStr, ok := sc.Keys["token"]; ok && len(tokenStr) > 0 {
			_, err := strconv.Atoi(tokenStr)
			if err == nil {
				return true
			}
		}
	}
	return false
}

// DoAuth uses the supplied username and pasword to authenticate with the
// SecurityCenter instance and stores an authentication token and session ID
// in the SC multiplexor to be used with future requests. Returns a boolean
// regarding the success or failure of the authentication attempt, along with
//an error if there was one.
func (sc *SC) DoAuth(username, password string) (bool, error) {
	res, err := sc.NewRequest("GET", "system").Do()
	if err != nil {
		return false, err
	}

	match := sessionRgx.FindStringSubmatch(res.HTTPResp.Header.Get("Set-Cookie"))
	if len(match) > 1 {
		sc.Keys["session"] = match[1]
	}

	res, err = sc.NewRequest("POST", "token", map[string]interface{}{
		"password": password,
		"username": username,
	}).Do()

	match = sessionRgx.FindStringSubmatch(res.HTTPResp.Header.Get("Set-Cookie"))
	if len(match) > 1 {
		sc.Keys["session"] = match[1]
	}

	if Verbose {
		var jsonData []byte
		jsonData, err = res.Data.MarshalJSON()
		fmt.Printf("headers: #%v\ngot auth data: %s\n", res.HTTPResp.Header, jsonData)
	}

	t := fmt.Sprint(res.Data.Get("response").Get("token").Interface())
	if len(t) > 0 {
		sc.Keys["token"] = t
		return true, nil
	}

	return false, nil
}
