package http

import (
	"fmt"
	"github.com/tkeel-io/tdtl"
	"net/url"

	"github.com/pkg/errors"
)

const (
	_pluginRudder      = "rudder"
	_pluginKeel        = "keel"
	_tenantExactMethod = "%s/apis/security/v1/tenants/exact"
	_tenantLoginMethod = "%s/apis/security/v1/oauth/%s/token"
)

func GetTenantID(host, tenantName string) (tenantID string, err error) {
	u, err := url.Parse(fmt.Sprintf(_tenantExactMethod, host))
	if err != nil {
		return "", errors.Wrap(err, "parse admin login method error")
	}
	val := u.Query()
	val.Set("title", tenantName)
	u.RawQuery = val.Encode()

	resp, err := Get(u.String())
	if err != nil {
		return "", errors.Wrap(err, "get tenant id error")
	}
	cc := tdtl.New(resp)
	tenantID = cc.Get("data.tenant_id").String()
	if tenantID != "" {
		return tenantID, nil
	}
	errMsg := cc.Get("msg").String()
	return "", errors.Wrap(fmt.Errorf(errMsg), "get tenant id error")
}

// /apis/security/v1/oauth/bfE2rmWK/token?grant_type=password&username=admin&password=123456
func GetTenantLoginToken(host, tenantID, username, password string) (accessToken string, refreshToken string, err error) {
	u, err := url.Parse(fmt.Sprintf(_tenantLoginMethod, host, tenantID))
	if err != nil {
		return "", "", errors.Wrap(err, "parse admin login method error")
	}
	val := u.Query()
	val.Set("grant_type", "password")
	val.Set("username", username)
	val.Set("password", password)
	u.RawQuery = val.Encode()

	resp, err := Get(u.String())
	if err != nil {
		return "", "", errors.Wrap(err, "invoking admin login err")
	}
	cc := tdtl.New(resp)
	accessToken = cc.Get("data.access_token").String()
	refreshToken = cc.Get("data.refresh_token").String()
	if accessToken != "" && refreshToken != "" {
		return accessToken, refreshToken, nil
	}
	errMsg := cc.Get("msg").String()
	return "", "", errors.Wrap(fmt.Errorf(errMsg), "get token error")
}

func Get(url string) (string, error) {
	fmt.Println(url)
	var (
		err error
		req *http.Request
	)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creat http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	var httpc http.Client

	var r *http.Response
	r, err = httpc.Do(req)
	if err != nil {
		return "", fmt.Errorf("error do http request: %w", err)
	}
	defer r.Body.Close()
	return readResponse(r)

}
func Post(url string, data []byte) (string, error) {
	fmt.Println(url)
	var (
		err error
		req *http.Request
	)
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("error creat http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

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
