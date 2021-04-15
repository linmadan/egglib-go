package comment

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"reflect"
)

func AddComments(db *pg.DB, model interface{}) {
	tableName := db.Model(model).TableModel().Table().SQLName
	columnsMap := make(map[string]*orm.Field)
	columns := db.Model(model).TableModel().Table().Fields
	for _, item := range columns {
		columnsMap[item.GoName] = item
	}
	valueOf := reflect.TypeOf(model)
	for i := 0; i < valueOf.Elem().NumField(); i++ {
		field := valueOf.Elem().Field(i)
		comment := field.Tag.Get("comment")
		if comment != "" {
			if field.Name == "tableName" {
				_, _ = db.Exec(fmt.Sprintf("COMMENT ON TABLE public.%s IS '%s';", tableName, comment))
			} else {
				if columnField, ok := columnsMap[field.Name]; ok {
					_, _ = db.Exec(fmt.Sprintf("COMMENT ON COLUMN public.%s.%s IS '%s';", tableName, columnField.SQLName, comment))
				}
			}
		}
	}
}
