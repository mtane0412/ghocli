/**
 * filter.go
 * フィールドフィルタリング機能
 *
 * 指定されたフィールドのみを出力する機能を提供します。
 */

package outfmt

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// FilterFields は指定されたフィールドのみをフィルタリングして出力します
//
// データがmap[string]interface{}またはそのスライスの場合、指定されたフィールドのみを抽出します。
// fieldsがnilまたは空の場合は、全フィールドをそのまま出力します。
func FilterFields(formatter *Formatter, data interface{}, fields []string) error {
	// フィールド指定がない場合は、そのまま出力
	if len(fields) == 0 {
		return formatter.Print(data)
	}

	// データの型に応じてフィルタリング
	switch v := data.(type) {
	case map[string]interface{}:
		// 単一のマップをフィルタリング
		filtered := filterMap(v, fields)
		return formatter.Print(filtered)

	case []map[string]interface{}:
		// スライスの各要素をフィルタリング
		filtered := make([]map[string]interface{}, len(v))
		for i, item := range v {
			filtered[i] = filterMap(item, fields)
		}

		// モードに応じて出力
		if formatter.mode == "plain" {
			return formatter.PrintTable(fields, mapSliceToRows(filtered, fields))
		}
		return formatter.Print(filtered)

	case []interface{}:
		// interface{}スライスをmap[string]interface{}スライスに変換
		var mapSlice []map[string]interface{}
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				mapSlice = append(mapSlice, m)
			}
		}
		if len(mapSlice) > 0 {
			filtered := make([]map[string]interface{}, len(mapSlice))
			for i, item := range mapSlice {
				filtered[i] = filterMap(item, fields)
			}

			// モードに応じて出力
			if formatter.mode == "plain" {
				return formatter.PrintTable(fields, mapSliceToRows(filtered, fields))
			}
			return formatter.Print(filtered)
		}

		// 変換できない場合はそのまま出力
		return formatter.Print(v)

	default:
		// 構造体の場合はreflectionを使用してフィルタリング
		return filterStruct(formatter, data, fields)
	}
}

// filterMap はマップから指定されたフィールドのみを抽出します
func filterMap(m map[string]interface{}, fields []string) map[string]interface{} {
	filtered := make(map[string]interface{})
	for _, field := range fields {
		if value, ok := m[field]; ok {
			filtered[field] = value
		}
	}
	return filtered
}

// mapSliceToRows はマップスライスをテーブル行に変換します
func mapSliceToRows(data []map[string]interface{}, fields []string) [][]string {
	rows := make([][]string, len(data))
	for i, item := range data {
		row := make([]string, len(fields))
		for j, field := range fields {
			if value, ok := item[field]; ok {
				row[j] = fmt.Sprintf("%v", value)
			} else {
				row[j] = ""
			}
		}
		rows[i] = row
	}
	return rows
}

// filterStruct は構造体から指定されたフィールドのみを抽出します
func filterStruct(formatter *Formatter, data interface{}, fields []string) error {
	// 構造体をマップに変換
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("構造体のマーシャルに失敗: %w", err)
	}

	var m interface{}
	if err := json.Unmarshal(jsonData, &m); err != nil {
		return fmt.Errorf("JSONのアンマーシャルに失敗: %w", err)
	}

	// マップに変換できた場合はフィルタリング
	switch v := m.(type) {
	case map[string]interface{}:
		filtered := filterMap(v, fields)
		return formatter.Print(filtered)
	case []interface{}:
		// スライスの場合
		var mapSlice []map[string]interface{}
		for _, item := range v {
			if itemMap, ok := item.(map[string]interface{}); ok {
				mapSlice = append(mapSlice, itemMap)
			}
		}
		if len(mapSlice) > 0 {
			filtered := make([]map[string]interface{}, len(mapSlice))
			for i, item := range mapSlice {
				filtered[i] = filterMap(item, fields)
			}

			// モードに応じて出力
			if formatter.mode == "plain" {
				return formatter.PrintTable(fields, mapSliceToRows(filtered, fields))
			}
			return formatter.Print(filtered)
		}
	}

	// フィルタリングできない場合はそのまま出力
	return formatter.Print(data)
}

// StructToMap は構造体をmap[string]interface{}に変換します
func StructToMap(data interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("data is not a struct")
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// JSONタグからフィールド名を取得
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// "omitempty"などのオプションを除去
		fieldName := strings.Split(jsonTag, ",")[0]
		if fieldName == "" {
			fieldName = field.Name
		}

		// 値が空の場合はスキップ（omitemptyの場合）
		if strings.Contains(jsonTag, "omitempty") && value.IsZero() {
			continue
		}

		result[fieldName] = value.Interface()
	}

	return result, nil
}
