package excel

type DataField struct {
	EnName string //字段英文名
	CnName string //字段中文名
}

type ExcelMaker interface {
	DataFieldList() []DataField                             //字段元数据列表
	CellValue(index int, enName string) (value interface{}) //获取单元格字段值
	DataListLen() int                                       //数据列表大小
	TableTitle() []string                                   //列表顶部自定义内容
}
