package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

const (
	LIMIT_PLAYER = 10 // ランキングの上限人数
)

func readInput() (Options, Entries, Scores, error) {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		return opts, nil, nil, err
	}

	entries, err := NewEntries(opts.Args.EntryFile)
	if err != nil {
		return opts, nil, nil, err
	}

	scores, err := NewScores(opts.Args.ScoreFile)
	if err != nil {
		return opts, nil, nil, err
	}

	return opts, entries, scores, nil
}

func processRanking(entries Entries, scores Scores) Rankings {
	rankings := NewRanking(entries, scores)
	rankingOptions := []RankingOption{WithSamePlayer(), WithDescendingScore(), WithSameRank(), WithSameRankPlayerID(), WithLimit(LIMIT_PLAYER)}
	rankings = applyRankingOptions(rankings, rankingOptions...)
	return rankings
}

func outputRanking(rankings Rankings) {
	rankings.Print()
}

func main() {
	_, entries, scores, err := readInput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	rankings := processRanking(entries, scores)

	outputRanking(rankings)
}
