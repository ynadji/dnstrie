package dnstrie

import (
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestCheckAndRemoveWildcard(t *testing.T) {
	type testCase struct {
		domain       string
		domainParsed string
		hasWildcard  bool
	}

	testCases := []testCase{
		testCase{"*.google.com", "google.com", true},
		testCase{"google.com", "google.com", false},
		testCase{"foo.*.google.com", "foo.*.google.com", false},
		testCase{"*google.com", "*google.com", false},
	}

	for _, tc := range testCases {
		parsed, wild := checkAndRemoveWildcard(tc.domain)
		if parsed != tc.domainParsed || wild != tc.hasWildcard {
			t.Fatalf("Failed with %+v. Got %v, %v.", tc, parsed, wild)
		}
	}
}

func TestReverseLabelSlice(t *testing.T) {
	type testCase struct {
		domain         string
		reversedLabels []string
	}

	testCases := []testCase{
		testCase{"www.google.com", []string{"com", "google", "www"}},
		testCase{"www.google.co.uk", []string{"uk", "co", "google", "www"}},
		testCase{"not.a.real.domain.asdashfkjah", nil},
		testCase{"not.a.real.!@#$.com", nil},
		testCase{"foo.com.gza.com", []string{"com", "gza", "com", "foo"}},
		testCase{"com", []string{"com"}},
		testCase{"", nil},
		testCase{"*.foo.com", []string{"com", "foo", "*"}},
	}

	for _, tc := range testCases {
		reversedLabels := reverseLabelSlice(tc.domain)
		if !reflect.DeepEqual(reversedLabels, tc.reversedLabels) {
			t.Fatalf("Failed to reverse labels. Got %+v expected %+v.", reversedLabels, tc.reversedLabels)
		}
	}
}
