package exceltool

import (
	"testing"
	"time"
)

type StructData struct {
	Id       int
	Name     string
	Ok       bool
	Number   float64
	CreateAt time.Time
}

type StructList []StructData

func (d StructList) FieldList() []Field {
	return []Field{
		{Key: "Name", Name: "描述"},
		{Key: "Id", Name: "序号"},
		{Key: "Ok", Name: "确认"},
		{Key: "Number", Name: "数量"},
		{Key: "CreateAt", Name: "时间"},
	}
}

func (d StructList) CellValue(index int, key string) (value interface{}) {
	switch key {
	case "Name":
		return d[index].Name
	case "Id":
		return d[index].Id
	case "Ok":
		return d[index].Ok
	case "Number":
		return d[index].Number
	case "CreateAt":
		return d[index].CreateAt.Format("2006/01/02 15:04")
	}
	return nil
}

func (d StructList) Len() int {
	return len(d)
}

func (d StructList) TableTitle() []string {
	return []string{
		"这是标题这是标题\n这是标题1",
		"这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2",
		"这是标题这是标题这是标题3",
	}
}

func TestMakeFromStruct(t *testing.T) {
	now := time.Now()
	demoData := []StructData{
		{1, "中文", false, 11.24, now},
		{2, "zhongwen", true, 0.3, now},
		{4, "中文", false, 11.24, now},
		{5, "zhongwen", true, 0.3, now},
	}
	exporter := NewExcelExport()
	err := exporter.ExportData(StructList(demoData), "txh")
	if err != nil {
		t.Error(err)
		return
	}
	err = exporter.SaveTo("fileDemo_struct.xlsx")
	if err != nil {
		t.Error(err)
		return
	}
}

type MapList []map[string]interface{}

func (d MapList) FieldList() []Field {
	return []Field{
		{Key: "Name", Name: "描述"},
		{Key: "Id", Name: "序号"},
		{Key: "Ok", Name: "确认"},
		{Key: "Number", Name: "数量"},
	}
}

func (d MapList) CellValue(index int, key string) (value interface{}) {
	if value, ok := d[index][key]; ok {
		return value
	}
	return nil
}

func (d MapList) Len() int {
	return len(d)
}

func (d MapList) TableTitle() []string {
	return nil
}

func TestMakeFromMap(t *testing.T) {
	demoData := []map[string]interface{}{
		{"Id": 1, "Name": "中文", "Number": 12.3, "Ok": true},
		{"Id": 2, "Name": "zhongwen", "Number": 11.3, "Ok": false},
	}
	exporter := NewExcelExport()
	err := exporter.ExportData(MapList(demoData), "")
	if err != nil {
		t.Error(err)
		return
	}
	err = exporter.SaveTo("fileDemo_map.xlsx")
	if err != nil {
		t.Error(err)
		return
	}
}
