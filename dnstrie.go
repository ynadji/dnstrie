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

// Match matches against exactly fully qualified domain names and zone
// wildcards.
func (root *DomainTrie) Match(domain string) bool {
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

func checkAndRemoveWildcard(domain string) (string, string) {
	if len(domain) < 2 {
		return domain, ""
	}
	if domain[1] == '.' && (domain[0] == '*' || domain[0] == '+') {
		return domain[2:], string(domain[0])
	}

	return domain, ""
}

func reverseLabelSlice(domain string) ([]string, error) {
	var reversedLabels []string
	domain, wildcard := checkAndRemoveWildcard(domain)
	labels := strings.Split(domain, ".")

	for i := len(labels) - 1; i >= 0; i-- {
		reversedLabels = append(reversedLabels, labels[i])
	}

	if wildcard != "" {
		reversedLabels = append(reversedLabels, wildcard)
	}

	return reversedLabels, nil
}

func addReversedLabelsToTrie(root *DomainTrie, reversedLabels []string) {
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

// MakeTrie returns the root of a trie given a slice of domain names.  Use
// dns.Normalize to prepare domains received from untrusted or unreliable
// sources.
func MakeTrie(domains []string) (*DomainTrie, error) {
	root := &DomainTrie{label: "."}

	for _, d := range domains {
		reversedLabels, err := reverseLabelSlice(d)
		length := len(reversedLabels)
		if err != nil {
			return nil, fmt.Errorf("Failed to build DomainTrie: %v", err)
		}
		// If it was plus, we need to add it both without the wildcard
		// for the exact match and with "*" for the normal wildcard
		// match.
		wasPlus := reversedLabels[length-1] == "+"
		if wasPlus {
			reversedLabels[length-1] = "*"
		}
		addReversedLabelsToTrie(root, reversedLabels)
		if wasPlus {
			addReversedLabelsToTrie(root, reversedLabels[:length-1])
		}
	}

	return root, nil
}
