package exceltool

import (
	"errors"
	"io"

	excelize "github.com/360EntSecGroup-Skylar/excelize/v2"
)

type Field struct {
	Key  string //字段英文名
	Name string //字段中文名
}

type ExcelListDataMaker interface {
	FieldList() []Field                                  //字段列表
	CellValue(index int, key string) (value interface{}) //获取字段值
	Len() int                                            //数据列总长度
	TableTitle() []string                                //列表顶部自定义内容
}

type ExcelExport struct {
	xlsx   *excelize.File
	sheets map[string]bool
}

func NewExcelExport() *ExcelExport {
	return &ExcelExport{
		xlsx: excelize.NewFile(),
		sheets: map[string]bool{
			"Sheet1": true,
		},
	}
}

//MakeListFromStruct 导出
func (export *ExcelExport) ExportData(sourData ExcelListDataMaker, newSheet string) error {
	//默认的表
	sheet := "Sheet1"
	if len(newSheet) > 0 {
		sheet = newSheet
		if _, ok := export.sheets[sheet]; !ok {
			export.xlsx.NewSheet(sheet)
			export.sheets[sheet] = true
		}
	}
	streamWriter, err := export.xlsx.NewStreamWriter(sheet)
	if err != nil {
		return err
	}
	fields := sourData.FieldList()
	if len(fields) == 0 {
		return errors.New("未设置数据表头")
	}
	titles := sourData.TableTitle()
	titleHeight := len(titles)
	titleWeight := len(fields)
	//设置顶部标题
	for i := range titles {
		begin, _ := excelize.CoordinatesToCellName(1, i+1) //n行第一列
		end, _ := excelize.CoordinatesToCellName(titleWeight, i+1)
		err := export.xlsx.MergeCell(sheet, begin, end)
		if err != nil {
			return err
		}
		streamWriter.SetRow(begin, []interface{}{titles[i]})
	}
	headerData := []interface{}{}
	//设置excel文档第titleHeight+1行的字段中文描述
	for index := range fields {
		cell := excelize.Cell{Value: fields[index].Name}
		headerData = append(headerData, cell)
	}
	//第titleHeight+1行A列。
	firstCell, err := excelize.CoordinatesToCellName(1, titleHeight+1)
	if err != nil {
		return err
	}
	streamWriter.SetRow(firstCell, headerData)

	//从excel第titleHeight+2行开始设置实际数据的值
	len := sourData.Len()
	for i := 0; i < len; i++ {
		var mainData []interface{}
		for _, field := range fields {
			cell := excelize.Cell{Value: sourData.CellValue(i, field.Key)}
			mainData = append(mainData, cell)
		}
		//A2，第titleHeight+2行A列开始
		beginCell, err := excelize.CoordinatesToCellName(1, i+titleHeight+2)
		if err != nil {
			return err
		}
		streamWriter.SetRow(beginCell, mainData)
	}
	return streamWriter.Flush()
}

func (export *ExcelExport) SaveTo(pathName string) error {
	return export.xlsx.SaveAs(pathName)
}

func (export *ExcelExport) WriteTo(w io.Writer) error {
	return export.xlsx.Write(w)
}
