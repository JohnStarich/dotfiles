package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	textrank2 "github.com/DavidBelicza/TextRank/v2"
	"github.com/JesusIslam/tldr"
)

func main() {
	args := args{
		Input:  os.Stdin,
		Output: os.Stdout,
	}
	flag.IntVar(&args.Sentences, "sentences", 5, "Number of sentences to summarize into.")
	flag.StringVar(&args.Strategy, "strategy", "textrank", "Number of sentences to summarize into. Available options: "+strings.Join(strategyNames(), ", "))
	flag.Parse()

	err := run(args)
	if err != nil {
		panic(err)
	}
}

type args struct {
	Input     io.Reader
	Output    io.Writer
	Sentences int
	Strategy  string
}

func run(args args) error {
	strategy, ok := strategies[args.Strategy]
	if !ok {
		fmt.Fprintf(args.Output, "Unknown strategy %q. Available options: %s\n", args.Strategy, strings.Join(strategyNames(), ", "))
		return nil
	}
	return strategy(args)
}

func formatSentence(s string) string {
	tokens := strings.FieldsFunc(s, unicode.IsSpace)
	return strings.Join(tokens, " ")
}

func readAll(r io.Reader) (string, error) {
	inputBytes, err := io.ReadAll(r)
	return string(inputBytes), err
}

var strategies = map[string]func(args) error{
	"textrank": textrank,
	"lexrank":  lexrank,
}

func strategyNames() []string {
	var names []string
	for name := range strategies {
		names = append(names, name)
	}
	return names
}

func textrank(args args) error {
	input, err := readAll(args.Input)
	if err != nil {
		return err
	}
	tr := textrank2.NewTextRank()
	rule := textrank2.NewDefaultRule()
	language := textrank2.NewDefaultLanguage()
	algorithmDef := textrank2.NewDefaultAlgorithm()
	tr.Populate(input, language, rule)
	tr.Ranking(algorithmDef)
	sentences := textrank2.FindSentencesByRelationWeight(tr, args.Sentences)
	for _, sentence := range sentences {
		fmt.Fprint(args.Output, formatSentence(sentence.Value), " ")
	}
	fmt.Fprintln(args.Output)
	return nil
}

func lexrank(args args) error {
	input, err := readAll(args.Input)
	if err != nil {
		return err
	}
	bag := tldr.New()
	result, err := bag.Summarize(input, args.Sentences)
	if err != nil {
		return err
	}
	for _, sentence := range result {
		fmt.Fprint(args.Output, formatSentence(sentence), " ")
	}
	fmt.Fprintln(args.Output)
	return nil
}
