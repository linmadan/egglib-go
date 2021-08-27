package excel

import (
	"github.com/xuri/excelize/v2"
	"io"
)

type ExcelImport struct {
	RowBegin    int    //第几行开始
	ColumnBegin int    //第几列开始，
	ColumnEnd   int    //第几列结束，
	Sheet       string //获取的表格
	DataFields  []DataField
}

func (excelImport ExcelImport) OpenExcelFromIoReader(r io.Reader) ([]map[string]string, error) {
	if excelImport.RowBegin < 1 {
		excelImport.RowBegin = 1
	}
	if excelImport.ColumnBegin < 1 {
		excelImport.ColumnBegin = 1
	}
	if excelImport.ColumnEnd == 0 {
		excelImport.ColumnEnd = len(excelImport.DataFields)
	}

	excelFile, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	rows, err := excelFile.Rows(excelImport.Sheet)
	if err != nil {
		return nil, err
	}
	var (
		rowDataList = make([]map[string]string, 0) //数据列表
		colHead     = make(map[int]string)         //map[索引数字]列表头映射英文字段
		rowIndex    int                            //行计数
	)
	for rows.Next() {
		rowIndex++
		cols, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		if rowIndex < excelImport.RowBegin {
			continue
		}
		if rowIndex == excelImport.RowBegin {
			//处理表头数据
			colHead = excelImport.transitionHead(cols)
			continue
		}
		rowData := make(map[string]string)

		for colIndex, colValue := range cols {
			if colIndex+1 < excelImport.ColumnBegin {
				continue
			}
			if colIndex > excelImport.ColumnEnd {
				break
			}
			if headEnName, ok := colHead[colIndex]; ok {
				rowData[headEnName] = colValue
			}
		}
		rowDataList = append(rowDataList, rowData)
	}
	return rowDataList, nil
}

func (excelImport ExcelImport) transitionHead(data []string) map[int]string {
	kv := make(map[int]string)
	for i := range data {
		for _, dataField := range excelImport.DataFields {
			if data[i] == dataField.CnName {
				kv[i] = dataField.EnName
			}
		}
	}
	return kv
}

func NewExcelImport() *ExcelImport {
	return &ExcelImport{
		RowBegin:    1,
		ColumnBegin: 1,
		Sheet:       "Sheet1",
	}
}
