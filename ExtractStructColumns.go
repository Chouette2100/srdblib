// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srdblib

import (
	"reflect"
	"strings"
)

const key = "json" // struct tag key

// ExtractStructColumns returns a comma-separated string of struct fields.
func ExtractStructColumns(model interface{}) string {
	columns := GetStructColumns(model)
	return strings.Join(columns, ", ")
}

// GetStructColumns returns a slice of struct fields.
func GetStructColumns(model interface{}) []string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	var columns []string
	collectColumns(t, &columns)
	return columns
}

// collectColumns collects struct fields recursively.
func collectColumns(t reflect.Type, columns *[]string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			// 埋め込みフィールドの場合、再帰的に処理
			collectColumns(field.Type, columns)
		} else {
			var columnName string
			if tag := field.Tag.Get(key); tag != "" {
				columnName = strings.Split(tag, ";")[0]
			} else {
				// jsonタグがない場合はフィールド名を小文字に変換して使用
				columnName = strings.ToLower(field.Name)
			}
			if columnName != "" {
				*columns = append(*columns, "`"+columnName+"`")
			}
		}
	}
}
