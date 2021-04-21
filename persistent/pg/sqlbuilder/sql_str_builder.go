package sqlbuilder

import (
	"strings"
)

func RemoveSqlFields(sqlBuildFields []string, strs ...string) []string {
	out := append([]string(nil), sqlBuildFields...)
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

func SqlFieldsSnippet(sqlBuildFields []string) string {
	return strings.Join(sqlBuildFields, ",")
}

func SqlPlaceHoldersSnippet(sqlBuildFields []string) string {
	var placeHolder []string
	for i := 0; i < len(sqlBuildFields); i++ {
		placeHolder = append(placeHolder, "?")
	}
	return strings.Join(placeHolder, ",")
}

func SqlUpdateFieldsSnippet(sqlBuildFields []string) string {
	return strings.Join(sqlBuildFields, "=?,") + "=?"
}
