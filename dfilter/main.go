package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/influxdata/influxdb/kit/cli"
	"github.com/ynadji/dnstrie"
)

var root *dnstrie.DomainTrie

var flags struct {
	matchFile string
	wildcard  bool
}

func readDomains(matchFilePath string) []string {
	f, err := os.Open(matchFilePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to read %s: %v", matchFilePath, err))
	}
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	return strings.Split(strings.TrimSpace(string(content)), "\n")
}

func run() error {
	// TODO: prob just switch to straight github.com/spf13/cobra
	if flags.matchFile == "" {
		fmt.Fprintln(os.Stderr, "Must provide match file!")
		os.Exit(2)
	}

	domains := readDomains(flags.matchFile)
	root, err := dnstrie.MakeTrie(domains)
	if err != nil {
		return fmt.Errorf("Failed to make trie: %v", err)
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		domain := scanner.Text()
		matched := root.Match(domain)

		if matched {
			fmt.Printf("%v\n", domain)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading standard input: %+v\n", err)
	}
	return nil
}

func main() {
	cmd := cli.NewCommand(&cli.Program{
		Run:  run,
		Name: "dfilter",
		Opts: []cli.Opt{
			{
				DestP: &flags.matchFile,
				Flag:  "matchFile",
				Desc:  "File of domain matches, one per line",
			},
		},
	})

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
