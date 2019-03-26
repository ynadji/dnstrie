// Package dnstrie creates a DNS-aware trie for fast filtering of domain names
// based on exact matches and wildcarded zone cut matches against a slice of
// matching criteria. For example, providing "*.org", "google.com", and
// "*.mail.google.com" would construct a trie like:
//
// tree
// +-- org
//     +-- *
// +-- com
//     +-- google
//         +-- mail
//             +-- *
// Where google.com and anything under but not including org, mail.google.com
// match the tree (with `tree.WildcardMatch`) and only google.com would match
// (with `tree.ExactMatch`).
package dnstrie

import (
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/ynadji/net/publicsuffix" // until changes are merged upstream
)

type DomainTrie struct {
	label  string
	others DomainTrieSlice
	end    bool
}

type DomainTrieSlice []*DomainTrie

// Empty returns true if nothing has been added to the trie and true otherwise.
func (root *DomainTrie) Empty() bool {
	return root.label == "" && root.others == nil && !root.end
}

// ExactMatch only matches against exactly fully qualified domain names and
// ignores zone wildcards.
func (root *DomainTrie) ExactMatch(domain string) bool {
	if !govalidator.IsDNSName(domain) || !publicsuffix.HasListedSuffix(domain) {
		return false
	}
	reversedLabels := reverseLabelSlice(domain)
	curr := root
	for _, label := range reversedLabels {
		node := findNode(label, curr.others)
		if node == nil {
			return false
		}
		curr = node
	}
	return curr.end
}

// WildcardMatch matches against exactly fully qualified domain names and zone
// wildcards.
func (root *DomainTrie) WildcardMatch(domain string) bool {
	if !govalidator.IsDNSName(domain) || !publicsuffix.HasListedSuffix(domain) {
		return false
	}
	reversedLabels := reverseLabelSlice(domain)
	curr := root
	for _, label := range reversedLabels {
		node := findNode("*", curr.others)
		if node != nil {
			return true
		}
		node = findNode(label, curr.others)
		if node == nil {
			return false
		}
		curr = node
	}
	return curr.end
}

func findNode(label string, others DomainTrieSlice) *DomainTrie {
	for _, trie := range others {
		if trie.label == label {
			return trie
		}
	}
	return nil
}

func checkAndRemoveWildcard(domain string) (string, bool) {
	if len(domain) < 2 {
		return domain, false
	}
	if domain[0] == '*' && domain[1] == '.' {
		return domain[2:], true
	} else {
		return domain, false
	}
}

func reverseLabelSlice(domain string) []string {
	var reversedLabels []string
	domain, wildcarded := checkAndRemoveWildcard(domain)
	if !govalidator.IsDNSName(domain) || !publicsuffix.HasListedSuffix(domain) {
		return nil
	}
	labels := strings.Split(domain, ".")

	for i := len(labels) - 1; i >= 0; i-- {
		reversedLabels = append(reversedLabels, labels[i])
	}

	if wildcarded {
		reversedLabels = append(reversedLabels, "*")
	}

	return reversedLabels
}

// MakeTrie returns the root of a trie given a slice of domain names. Invalid
// domains and those that do not use known TLDs are ignored.
func MakeTrie(domains []string) *DomainTrie {
	root := &DomainTrie{label: "."}

	for _, d := range domains {
		reversedLabels := reverseLabelSlice(d)
		if reversedLabels == nil {
			continue
		}

		curr := root
		for _, label := range reversedLabels {
			node := findNode(label, curr.others)
			if node == nil {
				node = &DomainTrie{label, DomainTrieSlice{}, false}
				curr.others = append(curr.others, node)
			}
			curr = node
		}
		curr.end = true
	}

	return root
}
