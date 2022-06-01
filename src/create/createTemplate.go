package create

import (
	"encoding/json"
	"errors"
	"fmt"
	"tkeelBatchTool/src/http"
	"tkeelBatchTool/src/parse"
)

var (
	event_count = "1"
)

func CreateTemplate(templateMap map[string]*parse.IotTemplate) error {
	for _, v := range templateMap {
		//create templateObj
		fmt.Println("start create template object\n")
		templateId, err := CreateTemplateObj(v.TemplateObj)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if templateId == "" {
			fmt.Println("templateID is empty")
			continue
		}
		fmt.Printf(templateId)

		//create properties
		fmt.Println("\nstart create attributes\n")
		err1 := createPropertie(templateId, v.Attributes, "attribute")
		if err1 != nil {
			fmt.Println(err1)
			return err1
		}
		fmt.Println("start create telemetry\n")
		err2 := createPropertie(templateId, v.Telemetry, "telemetry")
		if err2 != nil {
			fmt.Println(err2)
			return err2
		}
		fmt.Println("start create commands\n")
		err3 := createPropertie(templateId, v.Commands, "command")
		if err3 != nil {
			fmt.Println(err3)
			return err3
		}
	}
	return nil
}

func createPropertie(templateId string, propertie map[string]*parse.IotPropertie, clissify string) error {
	if len(propertie) == 0 {
        fmt.Println("point is empty")
		return nil
	}
	jsonstr, _ := json.Marshal(propertie)
	fmt.Printf("%s", string(jsonstr))
	resultMap, err := http.DoCreate(http.IotUrl, "/v1/templates/"+templateId+"/"+clissify, "POST", nil, jsonstr)
	if err != nil {
		return err
	}

	if resultMap["code"] != "io.tkeel.SUCCESS" {
		return fmt.Errorf(resultMap["msg"].(string))
	}
	return err
}

func CreateTemplateObj(temp parse.IotTemplateObj) (string, error) {
	jsonstr, _ := json.Marshal(temp)
	resultMap, err := http.DoCreate(http.IotUrl, "/v1/templates", "POST", nil, jsonstr)
	if err != nil {
		fmt.Println("create templateObj err\n")
		return "", err
	}

	code, ok := resultMap["code"]
	if ok != true {
		return "", errors.New("response err")
	}
	if code.(string) != "io.tkeel.SUCCESS" {
		return "", errors.New(resultMap["msg"].(string))
	}

	/*if templateObj, ok1 := resultMap["data"]; ok1 == true {
		if templateObj1, ok2 := templateObj.(map[string]interface{}); ok2 == true {
			if id, ok3 := templateObj1["id"]; ok3 == true {
				return id.(string), nil
			}
		}
	}*/ 
	return temp.Id, nil 
}
