package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

type PlayLog struct {
	CreatedAt  string
	PlayID     string
	Score      uint32
	HandleName string
}
type ranking struct {
	Rank       uint32
	PlayerID   string
	HandleName string
	Score      uint32
}

func ReadCSVFile(filepath string) ([][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func createRanking(playlogs map[string]PlayLog) ([]ranking, error) {
	// PlayLogからrankingへ変換
	var rlist []ranking
	for _, logItem := range playlogs {
		rlist = append(rlist, ranking{
			PlayerID:   logItem.PlayID,
			HandleName: logItem.HandleName,
			Score:      logItem.Score,
		})
	}
	// クイックソートでScoreの降順にソート
	quickSortRanking(rlist, 0, len(rlist)-1)
	// 順位番号の設定
	for i := range rlist {
		rlist[i].Rank = uint32(i + 1)
	}
	return rlist, nil
}

// クイックソートのヘルパー関数
func quickSortRanking(arr []ranking, low, high int) {
	if low < high {
		p := partition(arr, low, high)
		quickSortRanking(arr, low, p-1)
		quickSortRanking(arr, p+1, high)
	}
}

func partition(arr []ranking, low, high int) int {
	pivot := arr[high].Score
	i := low - 1
	for j := low; j < high; j++ {
		if arr[j].Score > pivot { // 降順
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

func main() {

	play_logs, err := ReadCSVFile("play_logs.csv")
	if err != nil {
		log.Fatal("Error reading play_logs.csv: %v", err)
		return

	}

	entries, err := ReadCSVFile("entries.csv")
	if err != nil {
		log.Fatal("Error reading entries.csv: %v", err)
		return
	}

	plalog := make(map[string]PlayLog, len(play_logs))

	for i, record := range play_logs {
		if i == 0 { // ヘッダー行をスキップ
			continue
		}
		plalog[record[1]] = PlayLog{
			CreatedAt: record[0],
			PlayID:    record[1],
			Score: func() uint32 {
				score, err := strconv.ParseUint(record[2], 10, 32)
				if err != nil {
					log.Fatalf("Error parsing score: %v", err)
				}
				return uint32(score)
			}(),
		}
	}
	for i, record := range entries {
		if i == 0 { // ヘッダー行をスキップ
			continue
		}
		entry := plalog[record[0]]
		entry.HandleName = record[1]
		plalog[record[0]] = entry
	}

	rankings, err := createRanking(plalog)
	if err != nil {
		log.Fatal("Error creating ranking: %v", err)
		return
	}
	for _, record := range rankings {
		fmt.Println(record.PlayerID)
		fmt.Println(record.HandleName)
		fmt.Println(record.Score)
		fmt.Println("\n")
	}

}
