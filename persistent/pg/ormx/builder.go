package ormx

import (
	"github.com/go-pg/pg/v10"
	"strings"
)

func FieldNames(model interface{}) []string {
	fieldIterator := pg.Model(model).TableModel().Table().Fields
	var fields []string
	for i := 0; i < len(fieldIterator); i++ {
		item := fieldIterator[i]
		if len(item.SQLName) == 0 {
			continue
		}
		fields = append(fields, string(item.SQLName))
	}
	return fields
}

func Remove(strings []string, strs ...string) []string {
	out := append([]string(nil), strings...)

	for _, str := range strs {
		var n int
		for _, v := range out {
			if v != str {
				out[n] = v
				n++
			}
		}
		out = out[:n]
	}

	return out
}

func ParamPlaceHolder(l int) string {
	var placeHolder []string
	for i := 0; i < l; i++ {
		placeHolder = append(placeHolder, "?")
	}
	return strings.Join(placeHolder, ",")
}
