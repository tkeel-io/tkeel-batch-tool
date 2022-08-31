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
	"time"
)

func CreateSpaceTree(spaceTreeMap map[string]*parse.SpaceNodeInfo, f *excelize.File, order []string) error {
	fmt.Println("start create spaceTree \n")
	for _, spaceNodeName := range order {
		node, okd := spaceTreeMap[spaceNodeName]
		if !okd {
			continue
		}
		time.Sleep(time.Duration(1000) * time.Millisecond)
		_, err := createSpaceTree(node)
		if err != nil {
			fmt.Println(err)
			continue
			//return err
		}
	}
	fmt.Println("end \n")
	return nil
}

func createSpaceTree(spaceNode *parse.SpaceNodeInfo) (string, error) {
	jsonstr, _ := json.Marshal(spaceNode)
	resultMap, err := http.DoCreate(conf.DefaultConfig.IotUrl, "/v1/groups", "POST", nil, jsonstr)
	if err != nil {
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
	return "", nil
}
