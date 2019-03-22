package dnstrie

import (
	"strings"

	"github.com/asaskevich/govalidator"
	"golang.org/x/net/publicsuffix"
)

type DomainTrie struct {
	label  string
	others DomainTrieSlice
	end    bool
}

type DomainTrieSlice []*DomainTrie

func (root *DomainTrie) Empty() bool {
	return root.label == "" && root.others == nil && !root.end
}

func (root *DomainTrie) ExactMatch(domain string) bool {
	if !govalidator.IsDNSName(domain) {
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

func (root *DomainTrie) WildcardMatch(domain string) bool {
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
	_, icann := publicsuffix.PublicSuffix(domain)
	if !govalidator.IsDNSName(domain) || !icann {
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
