package beego

import (
	"github.com/linmadan/egglib-go/utils/tool_funs"
	"reflect"
	"strings"
)

var supportTag = map[string]int{
	"_":         1,
	"path":      2,
	"recursion": 2,
}

func TransformDomainModelToBeegoModel(domainModel, beegoModel interface{}) error {
	beegoModelType := reflect.TypeOf(beegoModel).Elem()
	beegoModelValue := reflect.ValueOf(beegoModel).Elem()
	domainModelValue := reflect.ValueOf(domainModel).Elem()
	for i := 0; i < beegoModelType.NumField(); i++ {
		field := beegoModelType.Field(i)
		tagStr := field.Tag.Get("domain")
		attrs, tags, err := tool_funs.ParseStructTag(tagStr, supportTag)
		if err != nil {
			return err
		}
		noTransform := false
		for key := range attrs {
			switch key {
			case "_":
				noTransform = true
				break
			}
		}
		if noTransform {
			continue
		}
		for key, value := range tags {
			switch key {
			case "path":
				transformValue := domainModelValue
				for _, subPath := range strings.Split(value, `.`) {
					if transformValue.Kind() == reflect.Ptr {
						transformValue = transformValue.Elem()
					}
					transformValue = transformValue.FieldByName(subPath)
				}
				if field.Name == "Id" {
					switch reflect.TypeOf(transformValue.Interface()).Name() {
					case "string":
						beegoModelValue.FieldByName(field.Name).SetString(transformValue.Interface().(string))
					case "int":
						beegoModelValue.FieldByName(field.Name).SetInt(int64(transformValue.Interface().(int)))
					case "int64":
						beegoModelValue.FieldByName(field.Name).SetInt(transformValue.Interface().(int64))
					}
				} else {
					if beegoModelValue.FieldByName(field.Name) != reflect.ValueOf(nil) && transformValue != reflect.ValueOf(nil) {
						beegoModelValue.FieldByName(field.Name).Set(transformValue)
					}
				}
			case "recursion":
				if domainModelValue.FieldByName(value) != reflect.ValueOf(nil) && beegoModelValue.FieldByName(field.Name) != reflect.ValueOf(nil) {
					TransformDomainModelToBeegoModel(domainModelValue.FieldByName(value).Interface(), beegoModelValue.FieldByName(field.Name).Interface())
				}
			}
		}
	}
	return nil
}

func TransformBeegoModelToDomainModel(domainModel, beegoModel interface{}) error {
	beegoModelType := reflect.TypeOf(beegoModel).Elem()
	beegoModelValue := reflect.ValueOf(beegoModel).Elem()
	domainModelValue := reflect.ValueOf(domainModel).Elem()
	for i := 0; i < beegoModelType.NumField(); i++ {
		field := beegoModelType.Field(i)
		tagStr := field.Tag.Get("domain")
		attrs, tags, err := tool_funs.ParseStructTag(tagStr, supportTag)
		if err != nil {
			return err
		}
		noTransform := false
		for key := range attrs {
			switch key {
			case "_":
				noTransform = true
				break
			}
		}
		if noTransform {
			continue
		}
		for key, value := range tags {
			switch key {
			case "path":
				domainModelFieldValue := domainModelValue
				stopSet := false
				for _, subPath := range strings.Split(value, `.`) {
					if domainModelFieldValue.Kind() == reflect.Ptr {
						domainModelFieldValue = domainModelFieldValue.Elem()
					}
					if domainModelFieldValue == reflect.ValueOf(nil) {
						stopSet = true
						break
					}
					domainModelFieldValue = domainModelFieldValue.FieldByName(subPath)
				}
				if stopSet {
					continue
				}
				beegoModelFieldValue := beegoModelValue.FieldByName(field.Name)
				if domainModelFieldValue != reflect.ValueOf(nil) && beegoModelFieldValue != reflect.ValueOf(nil) {
					domainModelFieldValue.Set(beegoModelFieldValue)
				}
			case "recursion":
				if domainModelValue.FieldByName(value) != reflect.ValueOf(nil) && beegoModelValue.FieldByName(field.Name) != reflect.ValueOf(nil) {
					TransformBeegoModelToDomainModel(domainModelValue.FieldByName(value).Interface(), beegoModelValue.FieldByName(field.Name).Interface())
				}
			}
		}
	}
	return nil
}
