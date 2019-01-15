package sc

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"runtime"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

var (
	// SkipSSLVerify sets http.Client to skip SSL verification, useful for testing
	// in development environments without valid certificates
	SkipSSLVerify = false
	// TimeoutDuration is the default number of seconds to wait before requests to
	// SecurityCenter will timeout.
	TimeoutDuration = 90
	// Verbose truthfulness used to determine whether to log HTTP requests to Stdout
	Verbose     = false
	transporter *http.Transport
)

func initTransporter() {
	if transporter == nil {
		transporter = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: SkipSSLVerify,
			},
		}
	}
}

// NewRequest creates a new Request
func (sc *SC) NewRequest(method, path string, data ...map[string]interface{}) *Request {
	r := &Request{
		sc:     sc,
		method: method,
		path:   path,
	}

	if len(data) > 0 {
		r.data = mergeData(data...)
	} else {
		r.data = make(map[string]interface{})
	}

	return r
}

// NewSC returns a new instance of a SC request multiplexor
func NewSC(h string) *SC {
	return &SC{
		Host: h,
		Keys: map[string]string{
			"session": "",
			"token":   "",
		},
	}
}

// WithAuth authenticates the SC instance with the user-supplied
// session and tokens. Note this does nothing to validate the supplied
// session and token are valid.
func (sc *SC) WithAuth(session, token string) *SC {
	sc.Keys["session"] = session
	sc.Keys["token"] = token
	return sc
}

// Do does the http request prepared from Request r, returning
// a Response and an error if there was one
func (r *Request) Do() (*Response, error) {
	var (
		reqData  []byte
		jsonResp *simplejson.Json
	)

	initTransporter()

	uri, err := url.Parse(r.sc.Host)
	if err != nil {
		return nil, err
	}
	if len(uri.Scheme) == 0 {
		uri.Scheme = "https"
	}
	uri.Path += "/rest/" + r.path

	if r.method == "GET" {
		params := url.Values{}
		for k, v := range r.data {
			params.Add(k, fmt.Sprint(v))
		}
	} else if r.method == "POST" {
		reqData, err = json.Marshal(r.data)
		if err != nil {
			return nil, err
		}
	}

	client := &http.Client{
		Timeout:   time.Duration(time.Duration(TimeoutDuration) * time.Second),
		Transport: transporter,
	}

	req, err := http.NewRequest(r.method, uri.String(), bytes.NewBuffer(reqData))
	if err != nil {
		return nil, err
	}

	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])

	for key, val := range r.sc.Keys {
		switch key {
		case "session":
			if len(val) > 0 {
				req.AddCookie(&http.Cookie{
					Name:  "TNS_SESSIONID",
					Value: val,
				})
			}
		case "token":
			if len(val) > 0 && r.path != "token" && !strings.HasSuffix(f.Name(), "(*sc.SC).DoAuth") {
				req.Header.Add("X-SecurityCenter", val)
			}
		}
	}

	if r.method == "POST" || r.method == "PATCH" {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if Verbose {
		rawReq, _ := httputil.DumpRequestOut(req, true)
		rawResp, _ := httputil.DumpResponse(resp, true)

		fmt.Printf("req: %s\n========\nresp: %s\n========\n", string(rawReq[:]), string(rawResp[:]))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.Header.Get("Content-Type") == "application/json" {
		jsonResp, err = simplejson.NewJson(body)
	}

	return &Response{
		Status:   resp.StatusCode,
		URL:      resp.Request.URL.String(),
		Data:     jsonResp,
		RespBody: body,
		HTTPResp: resp,
	}, err
}
