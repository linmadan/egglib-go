package exceltool

// func TestReader(t *testing.T) {
// 	files, err := os.Open("./fileDemo_struct.xlsx")
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	defer files.Close()
// 	excelRd := NewExcelListReader()
// 	excelRd.ColumnBegin = 1
// 	excelRd.ColumnEnd = 5
// 	excelRd.RowBegin = 4
// 	excelRd.Sheet = "txh"
// 	excelRd.Fields = []Field{
// 		{Key: "name", Name: "描述"},
// 		{Key: "id", Name: "序号"},
// 		{Key: "ok", Name: "确认"},
// 		{Key: "number", Name: "数量"},
// 		{Key: "created", Name: "时间"},
// 	}
// 	data, err := excelRd.OpenReader(files)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	for i := range data {
// 		fmt.Printf("%+v", data[i])
// 	}
// }
