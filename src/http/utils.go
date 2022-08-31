/*
Copyright 2021 The tKeel Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"tkeelBatchTool/src/conf"
)

func Get(url string) (string, error) {
	var (
		err error
		req *http.Request
	)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creat http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+conf.DefaultConfig.Token)
	var httpc http.Client

	var r *http.Response
	r, err = httpc.Do(req)
	if err != nil {
		return "", fmt.Errorf("error do http request: %w", err)
	}
	if r.StatusCode < 200 || r.StatusCode >= 300 {
		ret, _ := readResponse(r)
		return "", fmt.Errorf("StatusCode: %d, error: %s", r.StatusCode, ret)
	}
	defer r.Body.Close()
	return readResponse(r)

}
func Post(url string, data []byte) (string, error) {
	var (
		err error
		req *http.Request
	)
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("error creat http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+conf.DefaultConfig.Token)

	var httpc http.Client

	var r *http.Response
	r, err = httpc.Do(req)
	if err != nil {
		return "", fmt.Errorf("error do http request: %w", err)
	}
	defer r.Body.Close()
	return readResponse(r)
}

func readResponse(response *http.Response) (string, error) {
	rb, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error read http response: %w", err)
	}

	if len(rb) > 0 {
		return string(rb), nil
	}

	return "", nil
}
