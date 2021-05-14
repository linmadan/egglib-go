package exceltool

import (
	"io"

	excelize "github.com/360EntSecGroup-Skylar/excelize/v2"
)

type ExcelListReader struct {
	RowBegin    int    //第几行开始
	ColumnBegin int    //第几列开始，
	ColumnEnd   int    //第几列结束，
	Sheet       string //获取的表格
	Fields      []Field
}

func NewExcelListReader() *ExcelListReader {
	return &ExcelListReader{
		RowBegin:    1,
		ColumnBegin: 1,
		Sheet:       "Sheet1",
	}
}

func (eRead ExcelListReader) OpenReader(r io.Reader) ([]map[string]string, error) {
	if eRead.RowBegin < 1 {
		eRead.RowBegin = 1
	}
	if eRead.ColumnBegin < 1 {
		eRead.ColumnBegin = 1
	}
	if eRead.ColumnEnd == 0 {
		eRead.ColumnEnd = len(eRead.Fields)
	}

	xlsxFile, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	rows, err := xlsxFile.Rows(eRead.Sheet)
	if err != nil {
		return nil, err
	}
	var (
		datas    = make([]map[string]string, 0) //数据列表
		listHead = make(map[int]string)         //map[索引数字]列表头映射英文字段
		rowIndex int                            //行计数
	)
	for rows.Next() {
		rowIndex++
		cols, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		if rowIndex < eRead.RowBegin {
			continue
		}
		if rowIndex == eRead.RowBegin {
			//处理表头数据
			listHead = eRead.transitionHead(cols)
			continue
		}
		rowData := make(map[string]string)

		for colK, colV := range cols {
			if colK > eRead.ColumnEnd {
				break
			}
			if colK+1 < eRead.ColumnBegin {
				continue
			}
			if headK, ok := listHead[colK]; ok {
				rowData[headK] = colV
			}
		}
		datas = append(datas, rowData)
	}
	return datas, nil
}

func (eRead ExcelListReader) transitionHead(data []string) map[int]string {
	kv := make(map[int]string)
	for i := range data {
		for _, field := range eRead.Fields {
			if data[i] == field.Name {
				kv[i] = field.Key
			}
		}
	}
	return kv
}
