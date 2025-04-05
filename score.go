package main

import (
	"fmt"
)

type Score struct {
	CreateTimestamp string `json:"create_timestamp"`
	PlayerID        string `json:"player_id"`
	Score           string `json:"score"`
}
type Scores []Score

func NewScores(filePath string) ([]Score, error) {
	// CSVから各構造体スライスに変換
	scores, err := UnmarshalCSV[Score](filePath)
	if err != nil {
		return nil, fmt.Errorf("エントリ読み込みエラー: %v", err)
	}

	scoreMap := []Score{}
	scoreMap = append(scoreMap, scores...)

	return scoreMap, nil
}
