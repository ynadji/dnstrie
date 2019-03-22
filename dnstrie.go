package dnstrie

import (
	"fmt"
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

func reverseLabelSlice(domain string) []string {
	var reversedLabels []string
	etld, icann := publicsuffix.PublicSuffix(domain)
	if !govalidator.IsDNSName(domain) || !icann {
		return nil
	}
	labels := strings.Split(domain, ".")
	foundTld := false
	tldParts := []string{}

	for i := len(labels) - 1; i >= 0; i-- {
		// This may be excessive. Surely people would want to just say
		// *.uk, right?  Well it was fun anyway. Unsure if it's worth
		// keeping this in now that I think about it. Go to bed.
		if !foundTld {
			tldParts = append(tldParts, labels[i])
			tldPartsCopy := make([]string, len(tldParts))
			copy(tldPartsCopy, tldParts)
			reverse(tldPartsCopy)
			possibleEtld := strings.Join(tldPartsCopy, ".")
			if possibleEtld == etld {
				foundTld = true
				labels[i] = possibleEtld
			} else {
				continue
			}
		}
		reversedLabels = append(reversedLabels, labels[i])
	}

	return reversedLabels
}

func buildTrie(domains []string) {
}
