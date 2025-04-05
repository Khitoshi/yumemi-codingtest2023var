package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

func readInput() (Options, Entries, Scores, error) {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		return opts, nil, nil, err
	}

	//entries := NewEntries("game_entry_log.csv")
	entries, err := NewEntries(opts.Args.EntryFile)
	if err != nil {
		return opts, nil, nil, err
	}

	//scores := NewScores("game_score_log.csv")
	scores, err := NewScores(opts.Args.ScoreFile)
	if err != nil {
		return opts, nil, nil, err
	}

	return opts, entries, scores, nil
}

func processRanking(entries Entries, scores Scores, opts Options) Rankings {
	rankings := NewRanking(entries, scores)
	rankingOptions := NewRankingOptions(opts)
	rankings = applyRankingOptions(rankings, rankingOptions...)
	return rankings
}

func outputRanking(rankings Rankings) {
	rankings.Print()
}

func main() {
	opts, entries, scores, err := readInput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	rankings := processRanking(entries, scores, opts)

	outputRanking(rankings)
}
