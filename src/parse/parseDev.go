package parse

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
	"tkeelBatchTool/src/conf"
)

//xlsx setting
var (
	devTableNameStartIndex = 1
	devStartRow            = 2

	spaceNameColNum = 1 //spacc name
	spaceIdColNum   = 2 //spacc Guid

	devNameColNum         = 3  //设备名
	devCustomIdColNum     = 4  //custom devid
	devDirectColNum       = 5  //直连
	devSelfLearnColNum    = 6  //子学习开关
	devTemplateNameColNum = 7  //模板名称
	devTemplateIdColNum   = 8  //模板ID
	devExtColNum          = 9  //扩展信息
	devDescColNum         = 10 //描述

	//devGuidColNum = 8 //dev guid  for del
)

//row meta data
type xlsxRowMetaDevData struct {
	spaceName string
	spaceId   string

	devName         string
	devCustomId     string
	devDirect       string
	devSelfLearn    string
	devTemplateName string
	devTemplateId   string
	devExt          string
	devDesc         string

	devDirectBool    bool
	devSelfLearnBool bool

	tableName string
	excelAxis string
	//devGuid string
	row int
}

// iot format propertie
type DevInfo struct {
	Name             string                 `json:"name"`
	CustomId         string                 `json:"customId"` //custom
	DirectConnection bool                   `json:"directConnection"`
	SelfLearn        bool                   `json:selfLearn`
	ParentName       string                 `json:"parentName"`
	ParentId         string                 `json:"parentId"`
	TemplateName     string                 `json:"templateName"`
	TemplateId       string                 `json:"templateId"`
	Extension        map[string]interface{} `json:"ext"`
	Description      string                 `json:"description"`

	//for retrography
	/*TableName string
	ExcelAxis string
	SpaceName string
	DelInfo   string
	Row       int*/
}

type Key struct {
	rowStart, rowEnd int
}

func createDevExt(xrmd xlsxRowMetaDevData) map[string]interface{} {
	//parse ext
	ext := make(map[string]interface{})

	json.Unmarshal([]byte(xrmd.devExt), &ext)
	return ext
}

func formatDevInfo(xrmd xlsxRowMetaDevData) (*DevInfo, error) {

	dev := &DevInfo{
		Name:        xrmd.devName,
		CustomId:    xrmd.devCustomId + "-" + conf.DefaultConfig.TenantId,
		Description: xrmd.devDesc,

		DirectConnection: xrmd.devDirectBool,
		SelfLearn:        xrmd.devSelfLearnBool,

		ParentName: xrmd.spaceName,
		ParentId:   xrmd.spaceId + "-" + conf.DefaultConfig.TenantId,

		TemplateName: xrmd.devTemplateName,
		TemplateId:   xrmd.devTemplateId + "-" + conf.DefaultConfig.TenantId,

		Extension: createDevExt(xrmd),

		//for
		//TableName: xrmd.tableName,
		//ExcelAxis: xrmd.excelAxis,
		//SpaceName: xrmd.spaceName,
		//SpaceCustomId:   xrmd.spaceCustomId,
		//SpaceDesc:       xrmd.spaceDesc,
		//DelInfo: xrmd.devGuid,
		//Row:     xrmd.row,
	}
	return dev, nil
}

func DoParseDevExcel(filePath string, sRow int, eRow int) (map[string]*DevInfo, *excelize.File, error, []string) {

	//container
	devMap := make(map[string]*DevInfo)
	order := make([]string, 10000)

	//open xlsx file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		fmt.Println(err)
		return devMap, nil, err, order
	}

	//get table list
	for index, tableName := range f.GetSheetMap() {
		fmt.Println(index, tableName)
		if index < devTableNameStartIndex {
			continue
		}

		//-------获取 table 上所有merge单元格 value-----------
		mergeMap := make(map[string]interface{})
		ma, err := f.GetMergeCells(tableName)
		if err != nil {
			fmt.Println(err)
			return devMap, nil, err, order
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
						vv.(map[Key]string)[Key{r, r1}] = value[1]
					} else {
						nv := make(map[Key]string)
						nv[Key{r, r1}] = value[1]
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
			return devMap, nil, err, order
		}

		for _, row := range rows {
			//find start row
			rowNum += 1
			if rowNum < devStartRow {
				continue
			}

			if sRow != 0 && rowNum < sRow {
				continue
			}
			if eRow != 0 && rowNum > eRow {
				continue
			}

			//parse row meta data
			var xrmd xlsxRowMetaDevData

			//----for  iot return id  \for del
			xrmd.tableName = tableName
			cStr, _ := excelize.ColumnNumberToName(devCustomIdColNum)
			axis, _ := excelize.JoinCellName(cStr, rowNum)
			xrmd.excelAxis = axis
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
						for k, v := range info.(map[Key]string) {
							if (rowNum >= k.rowStart) && (rowNum <= k.rowEnd) {
								colCell = v
								break
							}
						}
					}
				}
				//-------------------------------

				switch {
				case colNum == spaceNameColNum:
					xrmd.spaceName = strings.Trim(colCell, " ")
					break
				case colNum == spaceIdColNum:
					xrmd.spaceId = strings.Trim(colCell, " ")
					break

				case colNum == devNameColNum:
					xrmd.devName = strings.Trim(colCell, " ")
					break
				case colNum == devCustomIdColNum:
					xrmd.devCustomId = strings.Trim(colCell, " ")
					if xrmd.devCustomId == "" {
						xrmd.devCustomId = GetUUID()
						fmt.Println("devCustomId = ", xrmd.devCustomId)
						err := f.SetCellValue(xrmd.tableName, xrmd.excelAxis, xrmd.devCustomId)
						if err != nil {
							fmt.Println("write devUUID error")
							return nil, nil, errors.New("write devUUID error"), nil
						}
						f.Save()
					}
					break
				case colNum == devDirectColNum:
					xrmd.devDirect = strings.Trim(colCell, " ")
					break
				case colNum == devSelfLearnColNum:
					xrmd.devSelfLearn = strings.Trim(colCell, " ")
					if xrmd.devSelfLearn == "" {
						xrmd.devSelfLearn = "FALSE"
					}
					break
				case colNum == devTemplateNameColNum:
					xrmd.devTemplateName = strings.Trim(colCell, " ")
					break
				case colNum == devTemplateIdColNum:
					xrmd.devTemplateId = strings.Trim(colCell, " ")
					break
				case colNum == devExtColNum:
					xrmd.devExt = strings.Trim(colCell, " ")
					break
				case colNum == devDescColNum:
					xrmd.devDesc = strings.Trim(colCell, " ")
					break
				default:
					//fmt.Print("row parse error\n")
				}
				//fmt.Println(colCell)
			}

			//check excel value
			err := checkDevExcelValue(&xrmd)
			if err != nil {
				fmt.Println(err)
				return nil, nil, err, nil
			}

			info, _ := formatDevInfo(xrmd)
			if info != nil {
				devMap[xrmd.devName] = info
				order = append(order, xrmd.devName)
			}
		}
	}
	return devMap, f, nil, order
}

func checkDevExcelValue(xrmd *xlsxRowMetaDevData) error {
	//check value
	if xrmd.devName == "" {
		fmt.Println("row = ", xrmd.row)
		return errors.New("devName is null")
	}

	directbool, err := strconv.ParseBool(xrmd.devDirect)
	if err != nil {
		fmt.Println("row = ", xrmd.row)
		return errors.New("direct  is error")
	}
	xrmd.devDirectBool = directbool

	selfLearnbool, err1 := strconv.ParseBool(xrmd.devSelfLearn)
	if err1 != nil {
		fmt.Println("row = ", xrmd.row)
		return errors.New("SelfLearn  is error")
	}
	xrmd.devSelfLearnBool = selfLearnbool
	return nil

}
