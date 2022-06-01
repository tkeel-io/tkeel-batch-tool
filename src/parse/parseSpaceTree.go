package parse

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	//"strconv"
	"strings"
)

//xlsx setting
var (
	spaceTreeTableNameStartIndex = 1
	spaceTreeStartRow            = 2

	parentSpaceNodeNameColNum     = 1 //spacc name
	parentSpaceNodeCustomIdColNum = 2 //spacc Guid

	curSpaceNodeNameColNum     = 3 // 
	curSpaceNodeCustomIdColNum = 4 //

	curSpaceNodeExtColNum = 5

	curSpaceNodeDescColNum = 6
)

//row meta data
type xlsxRowMetaSpaceTreeData struct {
	parentSpaceNodeName     string //spacc name
	parentSpaceNodeCustomId string //spacc Guid
	curSpaceNodeName        string 
	curSpaceNodeCustomId    string 
	curSpaceNodeExt         string //ext
	curSpaceNodeDesc        string //

	tableName  string
	excelAxis  string
	excelAxis1 string

	row int
}

// iot format propertie
type SpaceNodeInfo struct {
	Name        string                 `json:"name"`
	CustomId    string                 `json:"customId"` //custom
	ParentName  string                 `json:"parentName"`
	ParentId    string                 `json:"parentId"`
	Description string                 `json:"description"`
	Extension   map[string]interface{} `json:"ext"`
}

type treeKey struct {
	rowStart, rowEnd int
}

func formatSpaceNodeInfo(xrmd xlsxRowMetaSpaceTreeData) (*SpaceNodeInfo, error) {

	spaceNode := &SpaceNodeInfo{
		Name:        xrmd.curSpaceNodeName,
		CustomId:    xrmd.curSpaceNodeCustomId,
		ParentName:  xrmd.parentSpaceNodeName,
		ParentId:    xrmd.parentSpaceNodeCustomId,
		Description: xrmd.curSpaceNodeDesc,
		Extension:   createSpaceNodeExt(xrmd),
	}
	return spaceNode, nil
}

func DoParseSpaceTreeExcel(filePath string, sRow int, eRow int) (map[string]*SpaceNodeInfo, *excelize.File, error,[]string) {

	//container
	spaceTreeMap := make(map[string]*SpaceNodeInfo)
	order := make([]string, 10000)
    uuidMap := make(map[string]string)

	//open xlsx file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		fmt.Println(err)
		return spaceTreeMap, nil, err,order 
	}

	//get table list
	for index, tableName := range f.GetSheetMap() {
		fmt.Println(index, tableName)
		if index < spaceTreeTableNameStartIndex {
			continue
		}

		//-------获取 table 上所有merge单元格 value-----------
		mergeMap := make(map[string]interface{})
		ma, err := f.GetMergeCells(tableName)
		if err != nil {
			fmt.Println(err)
			return spaceTreeMap, nil, err,order 
		}
		for _, value := range ma {
			if len(value) == 2 {
				str := strings.Split(value[0], " ")
				str1 := strings.Split(str[0], ":")
				c, r, _ := excelize.SplitCellName(str1[0])
				c1, r1, _ := excelize.SplitCellName(str1[1])
				if c == c1 {
					vv, ok := mergeMap[c]
					if ok {
						vv.(map[treeKey]string)[treeKey{r, r1}] = value[1]
					} else {
						nv := make(map[treeKey]string)
						nv[treeKey{r, r1}] = value[1]
						mergeMap[c] = nv
					}
				}
			}
		}
		//---------------------------------------------------

		// 获取 table 上所有单元格
		rowNum := 0
		rows, err := f.GetRows(tableName)
		if err != nil {
			fmt.Println(err)
			return spaceTreeMap, nil, err,order 
		}

		for _, row := range rows {
			//find start row
			rowNum += 1
			if rowNum < spaceTreeStartRow {
				continue
			}

			if sRow != 0 && rowNum < sRow {
				continue
			}
			if eRow != 0 && rowNum > eRow {
				continue
			}

			//parse row meta data
			var xrmd xlsxRowMetaSpaceTreeData

			//----for  iot return id  \for del
			xrmd.tableName = tableName
			cStr, _ := excelize.ColumnNumberToName(parentSpaceNodeCustomIdColNum)
			axis, _ := excelize.JoinCellName(cStr, rowNum)
			xrmd.excelAxis = axis

			cStr1, _ := excelize.ColumnNumberToName(curSpaceNodeCustomIdColNum)
			axis1, _ := excelize.JoinCellName(cStr1, rowNum)
			xrmd.excelAxis1 = axis1 
			xrmd.row = rowNum
			//-------------------------
			for colNum, colCell := range row {
				colNum += 1

				//-----获取 merge 单元格value----
				if colCell == "" {
					colStr, err := excelize.ColumnNumberToName(colNum)
					if err != nil {
						fmt.Println(err)
						continue
					}
					info, ok := mergeMap[colStr]
					if ok {
						for k, v := range info.(map[treeKey]string) {
							if (rowNum >= k.rowStart) && (rowNum <= k.rowEnd) {
								colCell = v
								break
							}
						}
					}
				}
				//-------------------------------
				switch {
				case colNum == parentSpaceNodeNameColNum:
					xrmd.parentSpaceNodeName = strings.Trim(colCell, " ")
					break
				case colNum == parentSpaceNodeCustomIdColNum:
					xrmd.parentSpaceNodeCustomId = strings.Trim(colCell, " ")
					if xrmd.parentSpaceNodeCustomId == "" && xrmd.parentSpaceNodeName != "" {
                        uuid, ok := uuidMap[xrmd.parentSpaceNodeName]
                        if !ok {
						    xrmd.parentSpaceNodeCustomId = GetUUID()
                            uuidMap[xrmd.parentSpaceNodeName] = xrmd.parentSpaceNodeCustomId
                        } else {
                            xrmd.parentSpaceNodeCustomId = uuid 
                        }
						fmt.Println("parent",xrmd.parentSpaceNodeName, xrmd.parentSpaceNodeCustomId)
						err := f.SetCellValue(xrmd.tableName, xrmd.excelAxis, xrmd.parentSpaceNodeCustomId)
						if err != nil {
							fmt.Println("write parentSpaceNodeUUID error")
							return nil, nil, errors.New("write parentSpaceNodeUUID error"),order 
						}
						f.Save()
					}
					break
				case colNum == curSpaceNodeNameColNum:
					xrmd.curSpaceNodeName = strings.Trim(colCell, " ")
					break
				case colNum == curSpaceNodeCustomIdColNum:
					xrmd.curSpaceNodeCustomId = strings.Trim(colCell, " ")
					if xrmd.curSpaceNodeCustomId == "" && xrmd.curSpaceNodeName !=""{
                        _, ok := uuidMap[xrmd.curSpaceNodeName]
                        if !ok {
						    xrmd.curSpaceNodeCustomId = GetUUID()
                            uuidMap[xrmd.curSpaceNodeName] = xrmd.curSpaceNodeCustomId
                        } else {
							return nil, nil, errors.New("curSpaceNodeName repeat error"),order
                        }
                        
						fmt.Println("cur",xrmd.curSpaceNodeName, xrmd.curSpaceNodeCustomId)
						err := f.SetCellValue(xrmd.tableName, xrmd.excelAxis1, xrmd.curSpaceNodeCustomId)
						if err != nil {
							fmt.Println("write curSpaceNodeUUID error")
							return nil, nil, errors.New("write curSpaceNodeUUID error"),order
						}
						f.Save()
					}
					break
				case colNum == curSpaceNodeExtColNum:
					xrmd.curSpaceNodeExt = strings.Trim(colCell, " ")
					break
				case colNum == curSpaceNodeDescColNum:
					xrmd.curSpaceNodeDesc = strings.Trim(colCell, " ")
					break
				default:
					//fmt.Print("row parse error\n")
				}
				//fmt.Println(colCell)
			}

			//check excel value
			err := checkSpaceTreeExcelValue(&xrmd)
			if err != nil {
				fmt.Println(err)
				return nil, nil, err,order 
			}

			info, _ := formatSpaceNodeInfo(xrmd)
			if info != nil {
				spaceTreeMap[xrmd.curSpaceNodeName] = info
				order = append(order, xrmd.curSpaceNodeName)
			}
		}
	}
	return spaceTreeMap, f, nil, order
}

func checkSpaceTreeExcelValue(xrmd *xlsxRowMetaSpaceTreeData) error {
	//check value
	if xrmd.curSpaceNodeName == "" {
		fmt.Println("row = ", xrmd.row)
		return errors.New("curSpaceNodeName is null")
	}
	return nil

}

func createSpaceNodeExt(xrmd xlsxRowMetaSpaceTreeData) map[string]interface{} {
	//parse ext
	ext := make(map[string]interface{})

	json.Unmarshal([]byte(xrmd.curSpaceNodeExt), &ext)
	return ext
}
