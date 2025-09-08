// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srdblib

import (
	"reflect"
	"strings"
)

// const key = "json" // struct tag key
const key = "db" // struct tag key

// ExtractStructColumns returns a comma-separated string of struct fields.
// gorpを使う場合はこの関数は必要ない、gorpに完全に移行するまでの暫定措置
func ExtractStructColumns(model interface{}) string {
	columns := GetStructColumns(model)
	clmlist := strings.Join(columns, ", ")
	sname := GetStructName(model)
	// 2025-09-06以降に追加したフィールドを除外する
	switch sname {
	case "Event", "Wevent":
		clmlist = strings.Replace(clmlist, ", `cmode`", "", -1)
	default:
	}
	return clmlist
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
				if tag == "-" {
					continue // "-" タグは無視
				}
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

// GetStructName extracts the name of the struct from an interface value.
// If the value is a pointer to a struct, it returns the name of the struct it points to.
// If the value is not a struct or a pointer to a struct, it returns an empty string.
func GetStructName(model interface{}) string {
	if model == nil {
		return ""
	}

	// reflect.TypeOf(model) はインターフェースの動的な型を返します。
	// 例えば &Event{} の場合、*main.Event 型を返します。
	t := reflect.TypeOf(model)

	// もし型がポインタであれば、Elem() を使ってポインタが指す先の型を取得します。
	// 例: *main.Event -> main.Event
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 取得した型が構造体であれば、Name() を使って構造体名を取得します。
	// 例: main.Event -> Event
	if t.Kind() == reflect.Struct {
		return t.Name()
	}

	// 構造体でもポインタでもない場合は空文字列を返します。
	return ""
}
