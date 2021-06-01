package exceltool

import (
	"testing"
	"time"
)

type StructData2 struct {
	Id       int
	Name     string
	Name0    string
	Name1    string
	Name2    string
	Name3    string
	Name4    string
	Name5    string
	Name6    string
	Name7    string
	Name8    string
	Name9    string
	Name10   string
	Ok       bool
	Number   float64
	CreateAt time.Time
}

type StructList2 []StructData2

func (d StructList2) FieldList() []Field {
	return []Field{
		{Key: "Name", Name: "描述"},
		{Key: "Name0", Name: "描述"},
		{Key: "Name1", Name: "描述"},
		{Key: "Name2", Name: "描述"},
		{Key: "Name3", Name: "描述"},
		{Key: "Name4", Name: "描述"},
		{Key: "Name5", Name: "描述"},
		{Key: "Name6", Name: "描述"},
		{Key: "Name7", Name: "描述"},
		{Key: "Name8", Name: "描述"},
		{Key: "Name9", Name: "描述"},
		{Key: "Name10", Name: "描述"},
		{Key: "Id", Name: "序号"},
		{Key: "Ok", Name: "确认"},
		{Key: "Number", Name: "数量"},
		{Key: "CreateAt", Name: "时间"},
	}
}

func (d StructList2) CellValue(index int, key string) (value interface{}) {
	switch key {
	case "Name":
		return d[index].Name
	case "Name0":
		return d[index].Name0
	case "Name1":
		return d[index].Name1
	case "Name2":
		return d[index].Name2
	case "Name3":
		return d[index].Name3
	case "Name4":
		return d[index].Name4
	case "Name5":
		return d[index].Name5
	case "Name6":
		return d[index].Name6
	case "Name7":
		return d[index].Name7
	case "Name8":
		return d[index].Name8
	case "Name9":
		return d[index].Name9
	case "Name10":
		return d[index].Name10
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

func (d StructList2) Len() int {
	return len(d)
}

func (d StructList2) TableTitle() []string {
	return []string{
		"这是标题这是标题\n这是标题1",
		"这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2",
		"这是标题这是标题这是标题3",
	}
}

func BenchmarkMakeFromStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		now := time.Now()
		demoData := []StructData2{}
		for i := 0; i < 10000; i++ {
			temp := StructData2{
				Id:       1,
				Name:     "中文这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2",
				Name0:    "中文这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2",
				Name1:    "中文这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2",
				Name2:    "中文这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2",
				Name3:    "中文这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2",
				Name4:    "中文这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2",
				Name5:    "中文这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2",
				Name6:    "中文这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2",
				Name7:    "中文这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2",
				Name8:    "中文这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2",
				Name9:    "中文这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2",
				Name10:   "中文这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2这是标题这是标题这是标题2",
				Ok:       false,
				Number:   11.24,
				CreateAt: now,
			}
			demoData = append(demoData, temp)
		}
		exporter := NewExcelExport()
		err := exporter.ExportData(StructList2(demoData), "txh")
		if err != nil {
			b.Error(err)
			return
		}
		err = exporter.SaveTo("fileDemo_struct.xlsx")
		if err != nil {
			b.Error(err)
			return
		}
	}
}
