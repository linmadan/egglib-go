package utils

import (
	"fmt"
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
		m[relType.Field(i).Name] = elem.Field(i).Interface()
	}
	return m
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
