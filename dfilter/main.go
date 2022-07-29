package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/ynadji/dnstrie"
)

var root *dnstrie.DomainTrie

func readDomains(matchFilePath string) []string {
	f, err := os.Open(matchFilePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to read %s: %v", matchFilePath, err))
	}
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	return strings.Split(strings.TrimSpace(string(content)), "\n")
}

func run(c *cli.Context) error {
	domains := readDomains(c.String("matches"))
	root, err := dnstrie.MakeTrie(domains)
	if err != nil {
		return fmt.Errorf("Failed to make trie: %v", err)
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		domain := scanner.Text()
		matched := root.Match(domain)

		if matched && !c.Bool("complement") {
			fmt.Printf("%v\n", domain)
		} else if !matched && c.Bool("complement") {
			fmt.Printf("%v\n", domain)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading standard input: %+v\n", err)
	}
	return nil
}

func main() {
	app := &cli.App{
		Name:   "dfilter",
		Usage:  "cat domains.txt | dfilter ...",
		Action: run,
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "matches",
			Usage:    "Path to file of domain matches, one per line.",
			Required: true,
		},
		&cli.BoolFlag{
			Name:    "complement",
			Usage:   "Invert matches",
			Aliases: []string{"c"},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
