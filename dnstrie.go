package dnstrie

import (
	"strings"

	"github.com/asaskevich/govalidator"
	"golang.org/x/net/publicsuffix"
)

type DomainTrie struct {
	label  string
	others DomainTrieSlice
}

type DomainTrieSlice []DomainTrie

func reverse(tldPartsCopy []string) {
	for i := len(tldPartsCopy)/2 - 1; i >= 0; i-- {
		opp := len(tldPartsCopy) - 1 - i
		tldPartsCopy[i], tldPartsCopy[opp] = tldPartsCopy[opp], tldPartsCopy[i]
	}
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

func buildTrie(domains []string) {
}
