package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
)

func structToMap(s interface{}) map[string]string {
	m := make(map[string]string)
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" {
			tag = field.Name
		}
		m[tag] = fmt.Sprintf("%v", val.Field(i).Interface())
	}
	return m
}

// UnmarshalCSV はCSVファイルの1行目（ヘッダー）の項目名と構造体のJSONタグを照合し、各レコードを構造体スライスに変換するジェネリック関数
func UnmarshalCSV[T any](filePath string) ([]T, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("CSVファイルオープンエラー: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// ヘッダーの読み込み
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("CSVヘッダー読み込みエラー: %v", err)
	}
	// ヘッダー名と列インデックスのマッピング作成
	headerMap := make(map[string]int)
	for i, h := range header {
		headerMap[h] = i
	}

	var result []T
	var zero T
	elemType := reflect.TypeOf(zero)

	// CSVレコードごとに構造体を生成
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("CSV読み込みエラー: %v", err)
		}
		// 新たな構造体インスタンスを生成
		elem := reflect.New(elemType).Elem()
		// 各フィールドについてJSONタグと照合して値をセット
		for i := 0; i < elem.NumField(); i++ {
			fieldType := elemType.Field(i)
			jsonTag := fieldType.Tag.Get("json")
			idx, ok := headerMap[jsonTag]
			if !ok || idx >= len(record) {
				continue
			}
			elem.Field(i).SetString(record[idx])
		}
		result = append(result, elem.Interface().(T))
	}
	return result, nil
}
