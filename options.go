package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Options struct {
	Desc  bool `short:"d" long:"desc" description:"スコアを降順にソートする"`
	Limit int  `short:"l" long:"limit" default:"0" description:"上位ランキングの件数制限（0の場合は制限なし）"`

	Args struct {
		EntryFile string `positional-arg-name:"entry" required:"yes" description:"エントリCSVファイルパス"`
		ScoreFile string `positional-arg-name:"score" required:"yes" description:"スコアCSVファイルパス"`
	} `positional-args:"yes"`
}

// RankingOption はランキングスライスに対して任意の変換を行うフック関数
type RankingOption func(Rankings) Rankings

func NewRankingOptions(opts Options) []RankingOption {
	rankingOptions := []RankingOption{}
	if opts.Limit > 0 {
		rankingOptions = append(rankingOptions, WithLimit(opts.Limit))
	}

	if opts.Desc {
		rankingOptions = append(rankingOptions, WithDescendingScore())
	}
	return rankingOptions
}

// 上位n件に制限するオプション
func WithLimit(n int) RankingOption {
	return func(rankings Rankings) Rankings {
		if len(rankings) > n {
			return rankings[:n]
		}
		return rankings
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
