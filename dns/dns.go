// Package dns standardizes checks for domain validity. All functions support
// punycode and unicode domains (IDN) by default.
package dns

import (
	"fmt"
	"strings"

	"github.com/asaskevich/govalidator"
	"golang.org/x/net/idna"
	"golang.org/x/net/publicsuffix"
)

// Normalize is used to sanitize domain names so common inconsistencies do not
// occur. Domains flowing into the system from 3rd parties should always be
// normalized before storing and further processing. This function will:
// * trim trailing & leading whitespace
// * lowercase
// * convert to punycode
// An error is returned if the punycode conversion fails.
func Normalize(domain string) (string, error) {
	domain, err := idna.ToASCII(strings.ToLower(strings.TrimSpace(domain)))
	if err != nil {
		return domain, fmt.Errorf("Failed to normalize domain %s: %v", domain, err)
	}
	return domain, nil
}

// Valid returns true if the domain passes govalidator's DNS check and
// was parseable by IDNA.
func Valid(domain string) bool {
	domain, err := idna.ToASCII(strings.ToLower(strings.TrimSpace(domain)))
	if err != nil {
		return false
	}
	return govalidator.IsDNSName(domain)
}

// HasListedSuffix returns true if the domain has a TLD that appears on the
// public suffix list and false otherwise. Converts to ASCII to ensure suffix
// check succeeds but does no other normalization.
func HasListedSuffix(domain string) bool {
	domain, err := idna.ToASCII(domain)
	if err != nil {
		return false
	}
	ps, icann := publicsuffix.PublicSuffix(domain)
	// Only ICANN-managed domains can have a single label and
	// privately-managed domains must have multiple labels. If there is no
	// known suffix, `PublicSuffix` just returns the last label to `ps`
	// (e.g., single label). If it isn't managed by ICANN and does not
	// contain a '.', it must not be present on the list.
	return icann || (strings.IndexByte(ps, '.') >= 0)
}

// IsListedSuffix returns true iff the argument `etld` is on the public suffix
// list, e.g., is a public or private effective top-level domain (eTLD).
func IsListedSuffix(etld string) bool {
	etld, err := idna.ToASCII(etld)
	if err != nil {
		return false
	}
	ps, icann := publicsuffix.PublicSuffix(etld)
	return (icann || (strings.IndexByte(ps, '.') >= 0)) && ps == etld
}

// IsPossibleDomain returns true if `domain` syntactically matches RFC 1035 and
// is or uses a known public or private TLD. Converts to ASCII to ensure suffix
// check succeeds but does no other normalization.
func IsPossibleDomain(domain string) bool {
	domain, err := idna.ToASCII(domain)
	if err != nil {
		return false
	}
	return govalidator.IsDNSName(domain) && HasListedSuffix(domain)
}

// IsRegisterableDomain returns true if `domain` syntactically matches RFC 1035
// and uses a known public or private TLD but is not a TLD itself.  This means
// the domain could be registered and used by a 3rd party, but the check does
// not resolve to see if it has been registered.  Converts to ASCII to ensure
// suffix check succeeds but does no other normalization.
func IsRegisterableDomain(domain string) bool {
	domain, err := idna.ToASCII(domain)
	if err != nil {
		return false
	}
	return govalidator.IsDNSName(domain) && HasListedSuffix(domain) && !IsListedSuffix(domain)
}
