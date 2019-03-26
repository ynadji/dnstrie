# dnstrie

A simple trie for filtering DNS names. Includes a CLI interface in the binary
`dfilter`.

## Install

$ go get -u github.com/ynadji/dnstrie

## dfilter

`dfilter` is a simple CLI tool for filtering domain names. Given a `--matchFile`
of fully qualified domains and wildcarded zone cuts to match on, `dfilter` will
take in domains on `STDIN` and print those that match the filter to `STDOUT`.

### Install

$ go get -u github.com/ynadji/dnstrie/dfilter

### Usage
```
$ dfilter -h
Usage:
  dfilter [flags]

Flags:
  -h, --help           help for dfilter
  --matchFile string   File of domain matches, one per line
  --wildcard           Accept wildcard matches
```
