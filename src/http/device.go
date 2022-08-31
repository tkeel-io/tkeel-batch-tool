package http

import (
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/tkeel-io/tdtl"
	"math/rand"
	"net/url"
	"path/filepath"
	"text/template"
	"time"

	"github.com/pkg/errors"
)

const (
	_deviceMethod = "%s/apis/tkeel-device/v1/devices/%s"
)

func GetDevice(host, deviceID string) (deviceInfo string, err error) {
	u, err := url.Parse(fmt.Sprintf(_deviceMethod, host, deviceID))
	if err != nil {
		return "", errors.Wrap(err, "parse admin login method error")
	}

	resp, err := Get(u.String())
	if err != nil {
		return "", errors.Wrap(err, "get device info error")
	}
	return string(resp), errors.Wrap(err, "get device info error")
}



func DeviceDataTemplate(templatePath string) (*template.Template, error) {
	rand.Seed(time.Now().Unix())
	funcMap := template.FuncMap{
	}
	for name, fn := range sprig.FuncMap() {
		funcMap[name] = fn
	}
	return template.New(filepath.Base(templatePath)).Funcs(funcMap).ParseFiles(templatePath)
}

// host = "http://preview.tkeel.io:30080/"
func GenDeviceTemplate(host, deviceID, mode string) ([]byte, error) {
	resultMap, err := GetDevice(host, deviceID)
	if err != nil {
		return nil, err
	}
	device := tdtl.New(resultMap)

	path := ""
	switch mode {
	case "telemetry":
		path = "data.deviceObject.configs.telemetry.define.fields"
	case "attributes":
		path = "data.deviceObject.configs.attributes.define.fields"
	default:
		return nil, fmt.Errorf("unknown mode type")
	}
	ret := tdtl.New("{}")
	device.Get(path).Foreach(func(key []byte, value *tdtl.Collect) {
		valueKey := value.Get("id").String()
		valueType := value.Get("type").String()
		ret.Set(valueKey, tdtl.StringNode(mockData(valueType)))
	})
	return ret.Raw(), nil
}


func mockData(typ string) string {
	switch typ {
	case "string":
		return `"mock_{{randAlphaNum 6}}"`
	case "bool":
		return `true`
	default:
		return `{{randInt 0 10}}.{{randInt 10 99}}`
	}
}
