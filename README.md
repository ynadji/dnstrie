# dnstrie

A simple trie for filtering DNS names. Includes a CLI interface in the binary
`dfilter`.

## Install

```
$ go get -u github.com/ynadji/dnstrie
```

## dfilter

`dfilter` is a simple CLI tool for filtering domain names. Given a `--matchFile`
of fully qualified domains and wildcarded zone cuts to match on, `dfilter` will
take in domains on `STDIN` and print those that match the filter to `STDOUT`.

### Install

```
$ go get -u github.com/ynadji/dnstrie/dfilter
```

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

#### Examples
```
$ echo "eff.org
random.foo.org
notareal.domain.test
google.com
mail.google.com
mine.mail.google.com" \
| dfilter --matchFile <(echo -e "*.org\ngoogle.com\n*.mail.google.com") --wildcard
eff.org
random.foo.org
google.com
mine.mail.google.com

$ echo "eff.org
random.foo.org
notareal.domain.test
google.com
mail.google.com
mine.mail.google.com" \
| dfilter --matchFile <(echo -e "*.org\ngoogle.com\n*.mail.google.com")
google.com
```
