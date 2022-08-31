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

func DelSpaceTree(spaceTreeMap map[string]*parse.SpaceNodeInfo, order []string) error {
	fmt.Println("del spaceTree")
	var ids []string
	idsMap := make(map[string]interface{})
	for i := len(order) - 1; i >= 0; i-- {
		name := order[i]
		node, okd := spaceTreeMap[name]
		if !okd {
			continue
		}
		ids = append(ids, node.CustomId)
	}

	idsMap["ids"] = ids
	fmt.Println("start del spaceTree: \n", ids)
	err := delspaceTree(idsMap)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("del spaceTree succeed  \n")
	}

	return nil
}

func delspaceTree(idsMap map[string]interface{}) error {
	jsonstr, _ := json.Marshal(idsMap)
	resultMap, err := http.DoCreate(conf.DefaultConfig.IotUrl, "/v1/groups/delete", "POST", nil, jsonstr)
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
