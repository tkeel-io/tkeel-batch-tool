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
	buf, err := json.Marshal(DefaultConfig)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, buf, 555)
	if err != nil {
		return err
	}
	return nil
}
