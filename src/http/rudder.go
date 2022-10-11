package http

import (
	"fmt"
	"net/url"

	"github.com/tkeel-io/tdtl"

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
