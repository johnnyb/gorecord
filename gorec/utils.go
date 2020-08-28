package gorec

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// This is useful for a simple conversion of url.Values (i.e., map[string][]string) into a format that is assignable on this system.
func MakeFormAssignable(origin map[string][]string) map[string]interface{} {
	newmap := map[string]interface{}{}
	for k, v := range origin {
		if len(v) == 0 {
			newmap[k] = nil
		} else {
			newmap[k] = v[0]
		}
	}
	return newmap
}

// This is similar to `MakeFormAssignable` but useful if you don't know if your request will be URL-formatted or JSON-formatted
func MakeRequestAssignable(req *http.Request) map[string]interface{} {
	ctype := req.Header.Get("Content-Type")
	if ctype == "application/x-www-form-urlencoded" || ctype == "multipart/form-data" {
		req.ParseForm()
		return MakeFormAssignable(req.Form)
	} else {
		newmap := map[string]interface{}{}
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return newmap
		}
		_ = json.Unmarshal(data, &newmap)
		return newmap
	}
}
