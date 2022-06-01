package conf

import (
	"encoding/json"
	"io/ioutil"
	"tkeelBatchTool/src/http"
)

//config
type Config struct {
	NorthUrl        string `json:"northUrl"`
	IotUrl          string `json:"iotUrl"`
	Token           string `json:"token"`
	AccessKey       string `json:"accessKey"`
	Secretaccesskey string `json:"secretaccesskey"`
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
	c := &Config{}
	err := c.load(path)
	if err != nil {
		return err
	}
	//set
	http.NorthUrl = c.NorthUrl
	http.IotUrl = c.IotUrl
	http.Token = c.Token
	http.AccessKey = c.AccessKey
	http.SecretAccessKey = c.Secretaccesskey
	return nil
}
