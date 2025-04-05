package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Options struct {
	Args struct {
		EntryFile string `positional-arg-name:"entry" required:"yes" description:"エントリCSVファイルパス"`
		ScoreFile string `positional-arg-name:"score" required:"yes" description:"スコアCSVファイルパス"`
	} `positional-args:"yes"`
}

// RankingOption はランキングスライスに対して任意の変換を行うフック関数
type RankingOption func(Rankings) Rankings

// 上位n件に制限するオプション
func WithLimit(n int) RankingOption {
	return func(rankings Rankings) Rankings {
		newRankings := make(Rankings, 0, n)
		for _, r := range rankings {
			if r.Rank > n {
				break
			}

			newRankings = append(newRankings, r)
		}

		return newRankings
	}
}

// 下位n件に制限するオプション
func WithLimitBottom(n int) RankingOption {
	return func(rankings Rankings) Rankings {
		if len(rankings) > n {
			return rankings[len(rankings)-n:]
		}
		return rankings
	}
}

// スコアの降順にソートするオプション
func WithDescendingScore() RankingOption {
	return func(rankings Rankings) Rankings {
		sort.Slice(rankings, func(i, j int) bool {
			return parseScore(rankings[i].Score) > parseScore(rankings[j].Score)
		})
		return rankings
	}
}

func parseScore(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "スコア(%s)の変換エラー: %v\n", s, err)
		return 0
	}
	return n
}

// フックを適用する関数
func applyRankingOptions(rankings Rankings, opts ...RankingOption) Rankings {
	for _, opt := range opts {
		rankings = opt(rankings)
	}
	return rankings
}

// スコアが同じ場合Rankを同じにするオプション
func WithSameRank() RankingOption {
	return func(rankings Rankings) Rankings {
		if len(rankings) == 0 {
			return rankings
		}
		// 先頭のランクを1に設定
		rankings[0].Rank = 1
		// 先頭のスコアを取得
		prevScore := parseScore(rankings[0].Score)
		prevRank := 1
		for i := 1; i < len(rankings); i++ {
			currentScore := parseScore(rankings[i].Score)
			if currentScore == prevScore {
				// 前のスコアと同じなら同一ランク
				rankings[i].Rank = prevRank
			} else {
				// 異なる場合は順位=i+1として更新
				rankings[i].Rank = i + 1
				prevRank = i + 1
				prevScore = currentScore
			}
		}
		return rankings
	}
}

// 重複する順位の場合はplyerIDでソートするオプション
func WithSameRankPlayerID() RankingOption {
	return func(rankings Rankings) Rankings {
		if len(rankings) == 0 {
			return rankings
		}
		sort.Slice(rankings, func(i, j int) bool {
			if rankings[i].Rank != rankings[j].Rank {
				return rankings[i].Rank < rankings[j].Rank
			}
			return rankings[i].PlayerID < rankings[j].PlayerID
		})
		return rankings
	}
}

// 同一プレイヤーが複数エントリしている場合、スコアが高い方のみを残すオプション
func WithSamePlayer() RankingOption {
	return func(rankings Rankings) Rankings {
		best := make(map[string]Ranking)
		for _, r := range rankings {

			if _, ok := best[r.PlayerID]; !ok {
				best[r.PlayerID] = r
			}

			if parseScore(r.Score) > parseScore(best[r.PlayerID].Score) {
				best[r.PlayerID] = r
			}
		}
		filtered := make(Rankings, 0, len(best))
		for _, b := range best {
			filtered = append(filtered, b)
		}

		return filtered
	}
}
