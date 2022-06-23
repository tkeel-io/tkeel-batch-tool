package del

import (
	"fmt"
	"tkeelBatchTool/src/conf"

	"encoding/json"
	"errors"
	//"strings"
	"tkeelBatchTool/src/http"
	"tkeelBatchTool/src/parse"
	//"os"
)

func DelMapper(mapperMap map[string]([]*parse.Expression)) error {
	for k, v := range mapperMap {
	    var paths []string
        pathsMap := make(map[string]interface{})
		for _, p := range v {
			paths = append(paths, p.Path)
		}
        pathsMap["paths"] = paths 
		err := delMapper(k, pathsMap)
		if err != nil {
			fmt.Println("del devId = %s, error = %s \n", k, err)
		} else {
			fmt.Println("del devId= %s mapper  succeed  \n", k)
		}
	}
	return nil
}

func delMapper(devId string, pathsMap map[string]interface{}) error {

	jsonstr, _ := json.Marshal(pathsMap)
	resultMap, err := http.DoCreate(conf.DefaultConfig.IotUrl, "/v1/devices/"+devId+"/relation/delete", "POST", nil, jsonstr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	code, ok := resultMap["code"]
	if ok != true {
		fmt.Println("response error ", devId)
		return errors.New("response err")
	}
	if code.(string) != "io.tkeel.SUCCESS" {
		fmt.Println(resultMap["msg"].(string))
		return errors.New(resultMap["msg"].(string))
	}
	return nil
}
