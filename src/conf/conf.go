package conf

import (
	"encoding/json"
	"io/ioutil"
)

var DefaultConfig = &Config{}

//config
type Config struct {
	NorthUrl        string `json:"northUrl"`
	IotUrl          string `json:"iotUrl"`
	Token           string `json:"token"`
	AccessKey       string `json:"access_key"`
	Secretaccesskey string `json:"secretaccesskey"`
	RefreshToken    string `json:"refresh_token"`
	Host            string `json:"host"`
	Tenant          string `json:"tenant"`
	TenantID        string `json:"tenant_id"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Broker          string `json:"broker"`
	DeviceID        string `json:"device_id"`
	DeviceToken     string `json:"device_token"`
	Template        string `json:"template"`
	TemplateMode    string `json:"template_mode"`
}

func (c *Config) load(path string) error {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, c)
	if err != nil {
		return err
	}
	return nil
}

func InitConfig(path string) error {
	err := DefaultConfig.load(path)
	if err != nil {
		return err
	}
	return nil
}

func SaveConfig(path string) error {
	buf, err := json.MarshalIndent(DefaultConfig, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, buf, 555)
	if err != nil {
		return err
	}
	return nil
}
