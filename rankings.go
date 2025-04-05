package main

import (
	"fmt"
	"reflect"
	"strings"
)

type Ranking struct {
	Rank       int    `json:"rank"`
	PlayerID   string `json:"player_id"`
	HandleName string `json:"handle_name"`
	Score      string `json:"score"`
}
type Rankings []Ranking

func NewRanking(entries Entries, scores Scores) Rankings {
	var rankings Rankings
	for _, score := range scores {

		entity, ok := entries.FindByPlayerID(score.PlayerID)
		if !ok {
			continue
		}

		entryMap := structToMap(entity)
		scoreMap := structToMap(score)
		var ranking Ranking
		// ランキング構造体の各フィールドを JSON タグで対応する値に設定
		rankingVal := reflect.ValueOf(&ranking).Elem()
		rankingType := rankingVal.Type()
		for i := 0; i < rankingType.NumField(); i++ {
			field := rankingType.Field(i)
			tag := field.Tag.Get("json")
			// score のマップにあるならそちら、なければ entry のマップを参照
			if val, ok := scoreMap[tag]; ok {
				rankingVal.Field(i).SetString(val)
			} else if val, ok := entryMap[tag]; ok {
				rankingVal.Field(i).SetString(val)
			}
		}
		// playerID の保証（scoreMapやentryMapになければ playerID をセット）
		if ranking.PlayerID == "" {
			ranking.PlayerID = score.PlayerID
		}
		ranking.Rank = len(rankings) + 1

		rankings = append(rankings, ranking)
	}
	return rankings
}

func (r Rankings) Print() {
	// ヘッダーを生成
	rankingType := reflect.TypeOf(Ranking{})
	headers := make([]string, rankingType.NumField())
	for i := 0; i < rankingType.NumField(); i++ {
		field := rankingType.Field(i)
		header := field.Tag.Get("json")
		if header == "" {
			header = field.Name
		}
		headers[i] = header
	}
	fmt.Println(strings.Join(headers, ","))

	// 各ランキングの出力
	for _, ranking := range r {
		value := reflect.ValueOf(ranking)
		row := make([]string, value.NumField())
		for j := 0; j < value.NumField(); j++ {
			row[j] = fmt.Sprintf("%v", value.Field(j).Interface())
		}
		fmt.Println(strings.Join(row, ","))
	}
}
