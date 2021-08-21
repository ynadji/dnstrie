# dnstrie

A simple trie for filtering DNS names. Includes a CLI interface in the binary
`dfilter`.

## Install

```
$ go get -u github.com/ynadji/dnstrie
```

## dfilter

`dfilter` is a simple CLI tool for filtering domain names. Given a file of of
fully qualified domains and wildcarded zone cuts to match on, `dfilter` will
take in domains on `STDIN` and print those that match the filter to
`STDOUT`. Matches can be specified with a leading `*`, which includes the parent
domain, or with a `+`, which only includes children. See the Example below.

### Install

```
$ go get -u github.com/ynadji/dnstrie/dfilter
```

### Usage
```
¡ dfilter -h
NAME:
   dfilter - cat domains.txt | dfilter ...

USAGE:
   dfilter [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --matches value  Path to file of domain matches, one per line.
   --help, -h       show help (default: false)
```

#### Example
```
$ echo "eff.org
random.foo.org
notareal.domain.test
google.com
mail.google.com
mine.mail.google.com
web.google.com
foo.web.google.com" \
| dfilter --matches <(echo -e "+.org\ngoogle.com\n+.mail.google.com\n*.web.google.com")
eff.org
random.foo.org
google.com
mine.mail.google.com
web.google.com
foo.web.google.com
```
