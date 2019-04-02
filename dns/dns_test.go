package dns

import (
	"testing"
)

func TestHasListedSuffix(t *testing.T) {
	var hasListedSuffixTestCases = []struct {
		domain string
		want   bool
	}{
		{"foo.com", true},
		{"test", false},
		{"com", true},
		{"foo.test", false}, // Reserved TLDs see https://tools.ietf.org/html/rfc2606#page-2
		{"foo.example", false},
		{"foo.invalid", false},
		{"foo.localhost", false},
		{"example", false},
		{"invalid", false},
		{"localhost", false},
		{"foo.baaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaar", false}, // too long, can never be valid TLD
		{"万岁.中国", true},                        // Unicode
		{"xn--chqu66a.xn--fiqs8s", true},       // Above in punycode
		{"ésta.bien.es", true},                 // Unicode
		{"xn--sta-9la.bien.es", true},          // Above in punycode
		{"ياسين.الجزائر", true},                // Works with reverse directional unicode
		{"xn--mgby9cnc.xn--lgbbat1ad8j", true}, // Above in punycode (parts reversed)
		{"!@#$%^&*.com", true},                 // Does not check for invalid characters
		{"dyndns-at-work.com", true},
	}

	for _, tc := range hasListedSuffixTestCases {
		got := HasListedSuffix(tc.domain)
		if got != tc.want {
			t.Errorf("%q: got %v, want %v", tc.domain, got, tc.want)
		}
	}
}

func TestIsListedSuffix(t *testing.T) {
	var isListedSuffixTestCases = []struct {
		domain string
		want   bool
	}{
		{"foo.com", false},
		{"test", false},
		{"com", true},
		{"foo.test", false}, // Reserved TLDs see https://tools.ietf.org/html/rfc2606#page-2
		{"foo.example", false},
		{"foo.invalid", false},
		{"foo.localhost", false},
		{"example", false},
		{"invalid", false},
		{"localhost", false},
		{"foo.baaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaar", false}, // too long, can never be valid TLD
		{"万岁.中国", false},                  // Unicode
		{"中国", true},                      // Unicode
		{"xn--chqu66a.xn--fiqs8s", false}, // Above in punycode
		{"xn--fiqs8s", true},              // Above in punycode
		{"ésta.bien.es", false},           // Unicode
		{"es", true},                      // Unicode
		{"xn--sta-9la.bien.es", false},    // Above in punycode
		{"ياسين.كوم", false},              // Works with reverse directional unicode
		{"كوم", true},
		{"!@#$%^&*.com", false},
		{"dyndns-at-work.com", true},
	}

	for _, tc := range isListedSuffixTestCases {
		got := IsListedSuffix(tc.domain)
		if got != tc.want {
			t.Errorf("%q: got %v, want %v", tc.domain, got, tc.want)
		}
	}
}

func TestIsPossibleDomain(t *testing.T) {
	var isPossibleDomainTestCases = []struct {
		domain string
		want   bool
	}{
		{"foo.com", true},
		{"test", false},
		{"com", true},
		{"foo.test", false}, // Reserved TLDs see https://tools.ietf.org/html/rfc2606#page-2
		{"foo.example", false},
		{"foo.invalid", false},
		{"foo.localhost", false},
		{"example", false},
		{"invalid", false},
		{"localhost", false},
		{"foo.baaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaar", false}, // too long, can never be valid TLD
		{"万岁.中国", true},                  // Unicode
		{"xn--chqu66a.xn--fiqs8s", true}, // Above in punycode
		{"ésta.bien.es", true},           // Unicode
		{"xn--sta-9la.bien.es", true},    // Above in punycode
		{"ياسين.كوم", true},              // Works with reverse directional unicode
		{"xn--mgby9cnc.xn--fhbei", true}, // Above in punycode
		// Should fail with leading/trailing whitespace
		{"  google.com", false},
		{"google.com  ", false},
		// All of the above, but with invalid domain characters
		{"!foo.com", false},
		{"!test", false},
		{"!com", false},
		{"!foo.test", false},
		{"!foo.example", false},
		{"!foo.invalid", false},
		{"!foo.localhost", false},
		{"!example", false},
		{"!invalid", false},
		{"!localhost", false},
		{"!foo.baaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaar", false},
		{"!万岁.中国", false},
		{"!xn--chqu66a.xn--fiqs8s", false},
		{"¡ésta.bien!.es", false},
		{"¿xn--sta-9la.bien?.es", false},
		{"!ياسين.كوم", false},
		{"!xn--mgby9cnc.xn--fhbei", false},
		{"!@#$%^&*.com", false},
		{"dyndns-at-work.com", true},
	}

	for _, tc := range isPossibleDomainTestCases {
		got := IsPossibleDomain(tc.domain)
		if got != tc.want {
			t.Errorf("%q: got %v, want %v", tc.domain, got, tc.want)
		}
	}
}

func TestIsRegisterableDomain(t *testing.T) {
	var isRegisterableDomainTestCases = []struct {
		domain string
		want   bool
	}{
		{"foo.com", true},
		{"test", false},
		{"com", false},      // Rejects eTLDs
		{"foo.test", false}, // Reserved TLDs see https://tools.ietf.org/html/rfc2606#page-2
		{"foo.example", false},
		{"foo.invalid", false},
		{"foo.localhost", false},
		{"example", false},
		{"invalid", false},
		{"localhost", false},
		{"foo.baaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaar", false}, // too long, can never be valid TLD
		{"万岁.中国", true},                  // Unicode
		{"xn--chqu66a.xn--fiqs8s", true}, // Above in punycode
		{"ésta.bien.es", true},           // Unicode
		{"xn--sta-9la.bien.es", true},    // Above in punycode
		{"ياسين.كوم", true},              // Works with reverse directional unicode
		{"xn--mgby9cnc.xn--fhbei", true}, // Above in punycode
		// Should fail with leading/trailing whitespace
		{"  google.com", false},
		{"google.com  ", false},
		// All of the above, but with invalid domain characters
		{"!foo.com", false},
		{"!test", false},
		{"!com", false},
		{"!foo.test", false},
		{"!foo.example", false},
		{"!foo.invalid", false},
		{"!foo.localhost", false},
		{"!example", false},
		{"!invalid", false},
		{"!localhost", false},
		{"!foo.baaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaar", false},
		{"!万岁.中国", false},
		{"!xn--chqu66a.xn--fiqs8s", false},
		{"¡ésta.bien!.es", false},
		{"¿xn--sta-9la.bien?.es", false},
		{"!ياسين.كوم", false},
		{"!xn--mgby9cnc.xn--fhbei", false},
		{"!@#$%^&*.com", false},
		{"dyndns-at-work.com", false}, // Private eTLD
	}

	for _, tc := range isRegisterableDomainTestCases {
		got := IsRegisterableDomain(tc.domain)
		if got != tc.want {
			t.Errorf("%q: got %v, want %v", tc.domain, got, tc.want)
		}
	}
}

func TestNormalize(t *testing.T) {
	var normTestCases = []struct {
		unNormDomain string
		normDomain   string
	}{
		{"  google.com", "google.com"},
		{"google.com  ", "google.com"},
		{"GOOGLE.cOm  ", "google.com"},
		{"  GOOGLE.cOm", "google.com"},
		{"ياسين.كوم", "xn--mgby9cnc.xn--fhbei"},
		{"ésta.bien.es", "xn--sta-9la.bien.es"},
		{"万岁.中国", "xn--chqu66a.xn--fiqs8s"},
		{"google.com", "google.com"},
		{"*.google.com", "*.google.com"},
		{"*.gOOgle.cOm", "*.google.com"},
		{"   *.gOOgle.cOm ", "*.google.com"},
	}

	for _, tc := range normTestCases {
		got, err := Normalize(tc.unNormDomain)
		if err != nil {
			t.Errorf("%q: got %v, want %v (err: %v)", tc.unNormDomain, got, tc.normDomain, err)
		}
		if got != tc.normDomain {
			t.Errorf("%q: got %v, want %v", tc.unNormDomain, got, tc.normDomain)
		}
	}
}
