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
	"fmt"
	"strings"
)

// DomainTrie is a struct for the recursive DNS-aware trie data structure. The
// three members represent the current label ("." for the root), the list of
// children and if this label can be considered an ending state for the tree (to
// identify that there is an exact domain match at this point). This should not
// be used directly and should instead be created using `dnstrie.MakeTrie`.
type DomainTrie struct {
	label  string
	others domainTrieSlice
	end    bool
}

type domainTrieSlice []*DomainTrie

// Empty returns true if nothing has been added to the trie and true otherwise.
func (root *DomainTrie) Empty() bool {
	return root.others == nil && !root.end
}

// ExactMatch only matches against exactly fully qualified domain names and
// ignores zone wildcards.
func (root *DomainTrie) ExactMatch(domain string) bool {
	reversedLabels, err := reverseLabelSlice(domain)
	if err != nil {
		return false
	}
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
// wildcards. Note that `domain` _should not_ contain the '*' character. If the
// trie was constructed with wildcarded matches, this will accept them unlike
// `ExactMatch`.
func (root *DomainTrie) WildcardMatch(domain string) bool {
	reversedLabels, err := reverseLabelSlice(domain)
	if err != nil {
		return false
	}
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

func findNode(label string, others domainTrieSlice) *DomainTrie {
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
	}

	return domain, false
}

func reverseLabelSlice(domain string) ([]string, error) {
	var reversedLabels []string
	domain, wildcarded := checkAndRemoveWildcard(domain)
	labels := strings.Split(domain, ".")

	for i := len(labels) - 1; i >= 0; i-- {
		reversedLabels = append(reversedLabels, labels[i])
	}

	if wildcarded {
		reversedLabels = append(reversedLabels, "*")
	}

	return reversedLabels, nil
}

// MakeTrie returns the root of a trie given a slice of domain names.  Use
// dns.Normalize to prepare domains received from untrusted or unreliable
// sources.
func MakeTrie(domains []string) (*DomainTrie, error) {
	root := &DomainTrie{label: "."}

	for _, d := range domains {
		reversedLabels, err := reverseLabelSlice(d)
		if err != nil {
			return nil, fmt.Errorf("Failed to build DomainTrie: %v", err)
		}

		curr := root
		for _, label := range reversedLabels {
			node := findNode(label, curr.others)
			if node == nil {
				node = &DomainTrie{label, domainTrieSlice{}, false}
				curr.others = append(curr.others, node)
			}
			curr = node
		}
		curr.end = true
	}

	return root, nil
}
