package del

import (
	"encoding/json"
	"errors"
	"fmt"
	"tkeelBatchTool/src/conf"
	"tkeelBatchTool/src/http"
	"tkeelBatchTool/src/parse"
	//"os"
)

func DelTemplate(templateMap map[string]*parse.IotTemplate) error {
	fmt.Println("start del template object\n")
	var id []string
	for _, v := range templateMap {
		id = append(id, v.TemplateObj.Id)
	}
    if len(id) == 0 {
	    fmt.Println("del template id is null  ")
        return errors.New("del error") 
    }
    idsMap := make(map[string]interface{})
    idsMap["ids"] = id 
	err := delTemplate(idsMap)
	if err != nil {
		fmt.Println(err)
		fmt.Println("del template err ")
        return err 
	}
	fmt.Println("del template success")
	return nil
}

func delTemplate(ids map[string]interface{}) error {
	if len(ids) == 0 {
		return nil
	}
	jsonstr, _ := json.Marshal(ids)
	fmt.Printf("%s", string(jsonstr))
	resultMap, err := http.DoCreate(conf.DefaultConfig.IotUrl, "/v1/templates/delete", "POST", nil, jsonstr)
	if err != nil {
		fmt.Println("del templateObj err\n")
		return err
	}

	if resultMap["code"] != "io.tkeel.SUCCESS" {
		return errors.New(resultMap["msg"].(string))
	}
	return err
}
