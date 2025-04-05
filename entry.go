package main

import (
	"fmt"
)

type Entry struct {
	PlayerID   string `json:"player_id"`
	HandleName string `json:"handle_name"`
}
type Entries []Entry

func (entries Entries) FindByPlayerID(playerID string) (Entry, bool) {
	for _, e := range entries {
		if e.PlayerID == playerID {
			return e, true
		}
	}
	return Entry{}, false
}

func NewEntries(filePath string) ([]Entry, error) {
	// CSVから各構造体スライスに変換
	entries, err := UnmarshalCSV[Entry](filePath)
	if err != nil {
		return nil, fmt.Errorf("エントリ読み込みエラー: %v", err)
	}

	entryMap := []Entry{}
	entryMap = append(entryMap, entries...)

	return entryMap, nil
}
