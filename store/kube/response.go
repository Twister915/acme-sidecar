package kube

import "encoding/json"

type responseData struct {
	data []byte
	statusCode int
}

func (resp responseData) AsMap() (d map[string]interface{}, err error) {
	err = json.Unmarshal(resp.data, &d)
	return
}

func (resp responseData) Status() int {
	return resp.statusCode
}