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

func DelDev(devMap map[string]*parse.DevInfo) error {
	fmt.Println("del dev")
	var ids []string
    idsMap := make(map[string]interface{})
	for _, dev := range devMap {
		ids = append(ids, dev.CustomId)
	}

    idsMap["ids"] = ids 
	fmt.Println("start del dev: \n", ids)
	err := delDev(idsMap)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("del dev succeed  \n")
	}

	return nil
}

func delDev(idsMap map[string]interface{}) error {
	jsonstr, _ := json.Marshal(idsMap)
	resultMap, err := http.DoCreate(conf.DefaultConfig.IotUrl, "/v1/devices/delete", "POST", nil, jsonstr)
	if err != nil {
		return err
	}

	// todo
	code, ok := resultMap["code"]
	if ok != true {
		return errors.New("response err")
	}
	if code.(string) != "io.tkeel.SUCCESS" {
		return errors.New(resultMap["msg"].(string))
	}
	return nil
}
