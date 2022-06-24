package create

import (
	"fmt"
	"tkeelBatchTool/src/conf"
	"tkeelBatchTool/src/http"
	"tkeelBatchTool/src/parse"
	//"github.com/360EntSecGroup-Skylar/excelize/v2"
	"encoding/json"
	"errors"
	"github.com/xuri/excelize/v2"
	//"strings"
	//"os"
)

func CreateDev(devMap map[string]*parse.DevInfo, f *excelize.File, order []string) error {

	for _, devName := range order {
		dev, okd := devMap[devName]
		if !okd {
			continue
		}

		//create dev
		fmt.Println("start create dev \n")
		//dev.Group = devGroupId //move to dev group
		_, err := createDev(dev)
		if err != nil {
			fmt.Println(dev.Name, err)
			continue
			//return err
		}
	}
	return nil
}

func createDev(dev *parse.DevInfo) (string, error) {
	jsonstr, _ := json.Marshal(dev)
	//fmt.Printf("%s",string(jsonstr))
	resultMap, err := http.DoCreate(conf.DefaultConfig.IotUrl, "/v1/devices", "POST", nil, jsonstr)
	if err != nil {
		fmt.Println("1")
		return "", err
	}

	// todo
	code, ok := resultMap["code"]
	if ok != true {
		return "", errors.New("response err")
	}
	if code.(string) != "io.tkeel.SUCCESS" {
		return "", errors.New(resultMap["msg"].(string))
	}

	/*if devObj, ok1 := resultMap["data"]; ok1 == true {
		if dev1, ok2 := devObj.(map[string]interface{}); ok2 == true {
			if id, ok3 := dev1["id"]; ok3 == true {
				return id.(string), nil
			}
		}
	}*/
	return dev.CustomId, nil
}
