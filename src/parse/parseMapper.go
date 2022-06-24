package parse

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
	"tkeelBatchTool/src/conf"
)

//xlsx setting
var (
	mapperTableNameStartIndex = 1
	mapperStartRow            = 3

	dev1NameColNum          = 1 //目标设备
	dev1IdColNum            = 2 //目标设备ID
	dev1PointClassifyColNum = 3 //目标设备测点类型
	dev1PointDataNameColNum = 4 //目标设备测点数据字段名

	mapperTypeColNum = 5 //Mapper

	dev2NameColNum          = 6 //来源设备
	dev2IdColNum            = 7 //目标设备ID
	dev2PointClassifyColNum = 8 //来源设备测点类型
	dev2PointDataNameColNum = 9 //来源设备测点数据字段名

)

//row meta data
type xlsxRowMetaMapperData struct {
	dev1Name          string
	dev1Id            string
	dev1PointClassify string
	dev1PointDataName string //目标设备测点数据字段名

	MapperType string //Mapper

	dev2Name          string //来源设备
	dev2Id            string
	dev2PointClassify string //来源设备测点类型
	dev2PointDataName string //来源设备测点数据字段名

	row int
}

// iot format propertie
type Expression struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Expression  string `json:"expression"`
	Description string `json:"description"`
}

func formatExpression(xrmd xlsxRowMetaMapperData) (*Expression, error) {

	xrmd.dev2Id += "-" + conf.DefaultConfig.TenantId
	expression := &Expression{
		Name:        "fromBatchTool_" + xrmd.dev1Name,
		Path:        xrmd.dev1PointClassify + "." + xrmd.dev1PointDataName,
		Expression:  xrmd.dev2Id + "." + xrmd.dev2PointClassify + "." + xrmd.dev2PointDataName,
		Description: xrmd.dev2Id + "=" + xrmd.dev2Name,
	}
	return expression, nil
}

func DoParseMapperExcel(filePath string, sRow int, eRow int) (map[string]([]*Expression), *excelize.File, error) {

	//container
	MapperMap := make(map[string]([]*Expression))

	//open xlsx file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		fmt.Println(err)
		return MapperMap, nil, err
	}

	//get table list
	for index, tableName := range f.GetSheetMap() {
		fmt.Println(index, tableName)
		if index < mapperTableNameStartIndex {
			continue
		}

		//-------获取 table 上所有merge单元格 value-----------
		mergeMap := make(map[string]interface{})
		ma, err := f.GetMergeCells(tableName)
		if err != nil {
			fmt.Println(err)
			return MapperMap, nil, err
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
			return MapperMap, nil, err
		}

		for _, row := range rows {
			//find start row
			rowNum += 1
			if rowNum < mapperStartRow {
				continue
			}

			if sRow != 0 && rowNum < sRow {
				continue
			}
			if eRow != 0 && rowNum > eRow {
				continue
			}

			//parse row meta data
			var xrmd xlsxRowMetaMapperData

			//----for  iot return id  \for del
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
				case colNum == dev1NameColNum:
					xrmd.dev1Name = strings.Trim(colCell, " ")
					break
				case colNum == dev1IdColNum:
					xrmd.dev1Id = strings.Trim(colCell, " ")
					break
				case colNum == dev1PointClassifyColNum:
					xrmd.dev1PointClassify = strings.Trim(colCell, " ")
					break
				case colNum == dev1PointDataNameColNum:
					xrmd.dev1PointDataName = strings.Trim(colCell, " ")
					break

				case colNum == mapperTypeColNum:
					xrmd.MapperType = strings.Trim(colCell, " ")
					break

				case colNum == dev2NameColNum:
					xrmd.dev2Name = strings.Trim(colCell, " ")
					break
				case colNum == dev2IdColNum:
					xrmd.dev2Id = strings.Trim(colCell, " ")
					break
				case colNum == dev2PointClassifyColNum:
					xrmd.dev2PointClassify = strings.Trim(colCell, " ")
					break
				case colNum == dev2PointDataNameColNum:
					xrmd.dev2PointDataName = strings.Trim(colCell, " ")
					break

				default:
					//fmt.Print("row parse error\n")
				}
				//fmt.Println(colCell)
			}

			//check excel
			err := checkMapperExcelValue(&xrmd)
			if err != nil {
				fmt.Println(err)
				return nil, nil, err
			}

			array, ok := MapperMap[xrmd.dev1Id+"-"+conf.DefaultConfig.TenantId]
			if !ok {
				array = make([]*Expression, 0)
				MapperMap[xrmd.dev1Id+"-"+conf.DefaultConfig.TenantId] = array
			}
			expression, _ := formatExpression(xrmd)
			if expression != nil {
				MapperMap[xrmd.dev1Id+"-"+conf.DefaultConfig.TenantId] = append(MapperMap[xrmd.dev1Id+"-"+conf.DefaultConfig.TenantId], expression)
			}
		}
	}
	return MapperMap, f, nil
}
func checkMapperExcelValue(xrmd *xlsxRowMetaMapperData) error {
	//check value
	if xrmd.dev1Name == "" || xrmd.dev1Id == "" || xrmd.dev1PointClassify == "" || xrmd.dev1PointDataName == "" {
		fmt.Println("row =  ", xrmd.row)
		return errors.New("dev1 data is error ")
	}
	if xrmd.dev2Name == "" || xrmd.dev2Id == "" || xrmd.dev2PointClassify == "" || xrmd.dev2PointDataName == "" {
		fmt.Println("row =  ", xrmd.row)
		return errors.New("dev2 data is error ")
	}
	return nil
}
