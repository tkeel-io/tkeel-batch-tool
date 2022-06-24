package parse

import (
	"fmt"
	//"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/google/uuid"
	"github.com/tkeel-io/kit/log"
	"github.com/xuri/excelize/v2"
	//"strconv"
	"strings"
	//"time"
	"encoding/json"
	"errors"
	"tkeelBatchTool/src/conf"
)

//xlsx setting
var (
	templateTableNameStartIndex = 1 //开始表
	templateStartRow            = 2 //开始行

	templateNameColNum   = 1 //模板对象名
	templateNameIdColNum = 2 //模板对象自定义id
	//templateDescColNum   = 3 //模板对象说明
	pointNameColNum     = 3 //测点名称
	pointNameIdColNum   = 4 //测点ID
	pointDataTypeColNum = 5 //数据类型
	pointTypeColNum     = 6 //测点类型
	pointDefineColNum   = 7 //测点定义集合
	pointDescColNum     = 8 //测点说明
)

//row meta data
type xlsxRowMetaData struct {
	templateName   string
	templateNameId string
	//templateDesc   string

	pointName     string
	pointNameId   string
	pointDataType string

	pointType   string
	pointDefine string
	pointDesc   string

	tableName string
	excelAxis string
	write     bool
	row       int
}

// iot format propertie
type IotPropertie struct {
	Name        string                 `json:"name"`
	DataType    string                 `json:"type"`
	Define      map[string]interface{} `json:"define"`
	Description string                 `json:"description"`
	Id          string                 `json:"id"`
	//PointType string `json:"schemaType"`
	//Ext map[string]interface{} `json:"ext"`

	//for retrography
	/*TableName string
		ExcelAxis string
	    Write bool
		Row       int*/
}

//fomat template
type IotTemplateObj struct {
	Name        string `json :"name"`
	Id          string `json:"customId"`
	Description string `json:"description"`
}

//iot format Template
type IotTemplate struct {
	TemplateObj IotTemplateObj

	Attributes map[string]*IotPropertie
	Telemetry  map[string]*IotPropertie
	Commands   map[string]*IotPropertie
}

type rangeKey struct {
	rowStart, rowEnd int
}

//Template map
var tkeelTemplateMap map[string]*IotTemplate

func DoParseTemplateExcel(filePath string) (map[string]*IotTemplate, *excelize.File, error) {

	//container
	tkeelTemplateMap = make(map[string]*IotTemplate)

	//open xlsx file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		fmt.Println(err)
		return tkeelTemplateMap, nil, err
	}

	//get table list
	for index, tableName := range f.GetSheetMap() {
		fmt.Println(index, tableName)
		if index < templateTableNameStartIndex {
			continue
		}

		//-------获取 table 上所有merge单元格 value-----------
		mergeMap := make(map[string]interface{})
		ma, err := f.GetMergeCells(tableName)
		if err != nil {
			fmt.Println(err)
			return tkeelTemplateMap, nil, err
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
						vv.(map[rangeKey]string)[rangeKey{r, r1}] = value[1]
					} else {
						nv := make(map[rangeKey]string)
						nv[rangeKey{r, r1}] = value[1]
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
			return tkeelTemplateMap, nil, err
		}

		for _, row := range rows {
			//find start row
			rowNum += 1
			if rowNum < templateStartRow {
				continue
			}

			//parse row meta data
			var xrmd xlsxRowMetaData

			xrmd.tableName = tableName
			cStr, _ := excelize.ColumnNumberToName(templateNameIdColNum)
			axis, _ := excelize.JoinCellName(cStr, rowNum)
			xrmd.excelAxis = axis
			xrmd.row = rowNum

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
						for k, v := range info.(map[rangeKey]string) {
							if (rowNum >= k.rowStart) && (rowNum <= k.rowEnd) {
								colCell = v
								break
							}
						}
					}
				}
				//-------------------------------

				switch {
				case colNum == templateNameColNum:
					xrmd.templateName = colCell
					break
				case colNum == templateNameIdColNum:
					xrmd.templateNameId = colCell
					break
				/*case colNum == templateDescColNum:
				xrmd.templateDesc = colCell
				break*/

				case colNum == pointNameColNum:
					xrmd.pointName = colCell
					break
				case colNum == pointNameIdColNum:
					xrmd.pointNameId = colCell
					break
				case colNum == pointDataTypeColNum:
					if colCell == "" {
						xrmd.pointDataType = "float"
					} else {
						xrmd.pointDataType = colCell
					}
					break

				case colNum == pointTypeColNum:
					xrmd.pointType = colCell
					break
				case colNum == pointDefineColNum:
					xrmd.pointDefine = colCell
					break
				case colNum == pointDescColNum:
					xrmd.pointDesc = colCell
					break
				default:
					//fmt.Print("row parse error\n")
				}
				fmt.Println(colCell)
			}

			//check
			//check excel value
			err := checkTemplateExcelValue(xrmd)
			if err != nil {
				fmt.Println(err)
				return nil, nil, err
			}

			//format
			_, ok := tkeelTemplateMap[xrmd.templateName]
			if !ok {
				tkeelTemplateMap[xrmd.templateName] = &IotTemplate{Attributes: make(map[string]*IotPropertie), Telemetry: make(map[string]*IotPropertie), Commands: make(map[string]*IotPropertie)}
				tkeelTemplateMap[xrmd.templateName].TemplateObj.Name = xrmd.templateName
				if xrmd.templateNameId == "" {
					xrmd.templateNameId = GetUUID()
					log.Debug("templateId = ", xrmd.templateNameId)
					xrmd.write = true
				}
				tkeelTemplateMap[xrmd.templateName].TemplateObj.Id = xrmd.templateNameId
				//tkeelTemplateMap[xrmd.templateName].TemplateObj.Description = xrmd.templateDesc
				tkeelTemplateMap[xrmd.templateName].TemplateObj.Description = "tkeelbatchtool add "
			}
			if xrmd.write || xrmd.templateNameId == "" {
				err := f.SetCellValue(xrmd.tableName, xrmd.excelAxis, tkeelTemplateMap[xrmd.templateName].TemplateObj.Id)
				if err != nil {
					fmt.Println("write templateUUID error")
				}
				f.Save()
			}

			// point
			if xrmd.pointNameId == "" {
				log.Error("error: pointId is empty")
				continue
			}

			switch {
			case xrmd.pointType == "attribute" || xrmd.pointType == "DI" || xrmd.pointType == "DO" || xrmd.pointType == "AO":
				tkeelTemplateMap[xrmd.templateName].Attributes[xrmd.pointNameId] = formatPropertie(xrmd)
				break
			case xrmd.pointType == "telemetry" || xrmd.pointType == "AI":
				tkeelTemplateMap[xrmd.templateName].Telemetry[xrmd.pointNameId] = formatPropertie(xrmd)
				break
			case xrmd.pointType == "command":
				tkeelTemplateMap[xrmd.templateName].Commands[xrmd.pointNameId] = formatPropertie(xrmd)
				break
			default:
			}
		}
	}
	for _, v := range tkeelTemplateMap {
		v.TemplateObj.Id += "-" + conf.DefaultConfig.TenantId
	}
	//log.Debug("all templates ", tkeelTemplateMap)
	return tkeelTemplateMap, f, err
}

func formatPropertie(xrmd xlsxRowMetaData) *IotPropertie {
	return &IotPropertie{Id: xrmd.pointNameId, Description: xrmd.pointDesc, Name: xrmd.pointName, Define: createDefine(xrmd), DataType: xrmd.pointDataType}
	//TableName:xrmd.tableName,ExcelAxis : xrmd.excelAxis,Row :xrmd.row, Write: xrmd.write}
}

func createDefine(xrmd xlsxRowMetaData) map[string]interface{} {
	define := make(map[string]interface{})

	json.Unmarshal([]byte(xrmd.pointDefine), &define)
	//custom
	/*str3 := strings.Split(xrmd.pointDefine, ";")
	for _, v := range str3 {
		str4 := strings.Split(v, "=")
		if len(str4) == 2 {
			t0 := strings.Trim(str4[0], " ")
			t00 := strings.Trim(t0, "\n")
			t1 := strings.Trim(str4[1], " ")
			t11 := strings.Trim(t1, "\n")
			define[t00] = t11
		}
	}*/
	return define
}

// generate uuid
func GetUUID() string {
	id := uuid.New()
	return "iotd-" + id.String()
}

func checkTemplateExcelValue(xrmd xlsxRowMetaData) error {
	fmt.Println("row = ", xrmd.row)
	if xrmd.templateName == "" {
		return errors.New("error: templateName is empty")
	}
	if xrmd.pointNameId == "" {
		return errors.New("error: pointDataId is empty")
	}
	if xrmd.pointName == "" {
		xrmd.pointName = xrmd.pointNameId
	}
	if xrmd.pointDataType == "" {
		return errors.New("error: pointDataType is empty")
	}
	if xrmd.pointType == "" {
		return errors.New("error: pointType is error")
	}
	return nil
}
