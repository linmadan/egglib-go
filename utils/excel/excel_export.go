package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

type ExcelExport struct {
	ExcelFile *excelize.File
	Sheets    map[string]bool
}

func (excelExport *ExcelExport) ExportData(excelMaker ExcelMaker, sheetName string) error {
	//默认的表
	sheet := "Sheet1"
	if len(sheetName) > 0 {
		sheet = sheetName
		if _, ok := excelExport.Sheets[sheet]; !ok {
			excelExport.ExcelFile.NewSheet(sheet)
			excelExport.Sheets[sheet] = true
		}
	}
	streamWriter, err := excelExport.ExcelFile.NewStreamWriter(sheet)
	if err != nil {
		return err
	}
	dataFieldList := excelMaker.DataFieldList()
	if len(dataFieldList) == 0 {
		return fmt.Errorf("未设置数据表头")
	}
	titles := excelMaker.TableTitle()
	titleRows := len(titles)
	titleCols := len(dataFieldList)
	//设置顶部标题
	for i := range titles {
		begin, _ := excelize.CoordinatesToCellName(1, i+1) //n行第一列
		end, _ := excelize.CoordinatesToCellName(titleCols, i+1)
		err := excelExport.ExcelFile.MergeCell(sheet, begin, end)
		if err != nil {
			return err
		}
		streamWriter.SetRow(begin, []interface{}{titles[i]})
	}
	headerData := make([]interface{}, 0)
	//设置excel文档第titleRows+1行的字段中文描述
	for index := range dataFieldList {
		cell := excelize.Cell{Value: dataFieldList[index].CnName}
		headerData = append(headerData, cell)
	}
	//第titleRows+1行A列。
	firstCell, err := excelize.CoordinatesToCellName(1, titleRows+1)
	if err != nil {
		return err
	}
	streamWriter.SetRow(firstCell, headerData)

	//从excel第titleRows+2行开始设置实际数据的值
	len := excelMaker.DataListLen()
	for i := 0; i < len; i++ {
		var mainData []interface{}
		for _, dataField := range dataFieldList {
			cell := excelize.Cell{Value: excelMaker.CellValue(i, dataField.EnName)}
			mainData = append(mainData, cell)
		}
		//A2，第titleRows+2行A列开始
		beginCell, err := excelize.CoordinatesToCellName(1, i+titleRows+2)
		if err != nil {
			return err
		}
		streamWriter.SetRow(beginCell, mainData)
	}
	return streamWriter.Flush()
}

func NewExcelExport() *ExcelExport {
	return &ExcelExport{
		ExcelFile: excelize.NewFile(),
		Sheets: map[string]bool{
			"Sheet1": true,
		},
	}
}
