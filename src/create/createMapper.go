package create

import (
	"tkeelBatchTool/src/conf"
	"tkeelBatchTool/src/http"
	"tkeelBatchTool/src/parse"
	//"github.com/360EntSecGroup-Skylar/excelize/v2"
	"encoding/json"
	"errors"
	"github.com/tkeel-io/kit/log"
	"github.com/xuri/excelize/v2"
	//"strings"
	//"os"
)

func CreateMapper(mapperMap map[string]([]*parse.Expression), f *excelize.File) error {
	for k, v := range mapperMap {
		expMap := make(map[string]interface{})
		expMap["expressions"] = v
		jsonstr, _ := json.Marshal(expMap)
		resultMap, err := http.DoCreate(conf.DefaultConfig.IotUrl, "/v1/devices/"+k+"/relation", "POST", nil, jsonstr)
		if err != nil {
			log.Error(err)
			return err
		}

		code, ok := resultMap["code"]
		if ok != true {
			log.Error("response error ", k)
			return errors.New("response err")
		}
		if code.(string) != "io.tkeel.SUCCESS" {
			log.Error(resultMap["msg"].(string), k)
			//return errors.New(resultMap["msg"].(string))
			continue
		}

	}
	return nil
}
