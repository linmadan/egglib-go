package tool_funs

import (
	"fmt"
	"github.com/linmadan/egglib-go/utils/string_convert"
	"reflect"
	"strings"
)

func ParseStructTag(tagStr string, supportTag map[string]int) (map[string]bool, map[string]string, error) {
	attrs := make(map[string]bool)
	tags := make(map[string]string)
	for _, v := range strings.Split(tagStr, `;`) {
		if v == "" {
			continue
		}
		v = strings.TrimSpace(v)
		if t := strings.ToLower(v); supportTag[t] == 1 {
			attrs[t] = true
		} else if i := strings.Index(v, "("); i > 0 && strings.Index(v, ")") == len(v)-1 {
			name := t[:i]
			if supportTag[name] == 2 {
				v = v[i+1 : len(v)-1]
				tags[name] = v
			}
		} else {
			return attrs, tags, fmt.Errorf("不支持的转化标签")
		}
	}
	return attrs, tags, nil
}

func SimpleStructToMap(toMapStruct interface{}) map[string]interface{} {
	elem := reflect.ValueOf(toMapStruct).Elem()
	relType := elem.Type()
	m := make(map[string]interface{})
	for i := 0; i < relType.NumField(); i++ {
		m[string_convert.CamelCase(relType.Field(i).Name, false, false)] = elem.Field(i).Interface()
	}
	if pageNumber, ok := m["pageNumber"]; ok {
		var pageSize int64
		if _, ok := m["pageSize"]; ok {
			pageSize = m["pageSize"].(int64)
		} else {
			pageSize = 20
		}
		m["offset"] = (pageNumber.(int64) - 1) * pageSize
		m["limit"] = pageSize
	}
	return m
}

func SimpleWrapGridMap(total int64, list interface{}) map[string]interface{} {
	grid := map[string]interface{}{"total": total, "list": list}
	return map[string]interface{}{
		"grid": grid,
	}
}

func QueryOptionsStructToMap(toMapQueryOptionsStruct interface{}) map[string]interface{} {
	elem := reflect.ValueOf(toMapQueryOptionsStruct).Elem()
	relType := elem.Type()
	m := make(map[string]interface{})
	for i := 0; i < relType.NumField(); i++ {
		if relType.Field(i).Name == "BaseQueryOptions" {
			for j := 0; j < elem.Field(i).Type().NumField(); j++ {
				m[elem.Field(i).Type().Field(j).Name] = elem.Field(i).Field(j).Interface()
			}
		} else {
			m[relType.Field(i).Name] = elem.Field(i).Interface()
		}
	}
	return m
}
