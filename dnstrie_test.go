package dnstrie

import (
	"os"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
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

func TestMakeTrie(t *testing.T) {
	type testCase struct {
		domains []string
		root    *DomainTrie
	}

	testCases := []testCase{
		{[]string{"www.google.com", "*.google.com"}, &DomainTrie{
			label: ".",
			others: DomainTrieSlice{
				&DomainTrie{
					label: "com",
					others: DomainTrieSlice{
						&DomainTrie{
							label: "google",
							others: DomainTrieSlice{
								&DomainTrie{"www", DomainTrieSlice{}, true},
								&DomainTrie{"*", DomainTrieSlice{}, true},
							},
						},
					},
				},
			},
		},
		},
	}

	for _, tc := range testCases {
		root := MakeTrie(tc.domains)
		if !reflect.DeepEqual(root, tc.root) {
			t.Fatalf("Failed to MakeTrie. Got:\n%+v\nExpected:\n%+v\n", spew.Sdump(root), spew.Sdump(tc.root))
		}
	}
}

func TestExactMatch(t *testing.T) {
	type testCase struct {
		domain string
		match  bool
	}
	root := MakeTrie([]string{"*.google.com", "www.google.org", "*.biz", "notarealdomain", "*nadji.us", "onizuka.homelinux.org"})

	testCases := []testCase{
		testCase{"www.google.org", true},
		testCase{"www.google.com", false},
		testCase{"google.com", false},
		testCase{"google.biz", false},
		testCase{"foo.google.biz", false},
		testCase{"bar.foo.google.biz", false},
		testCase{"notarealdomain", false},
		testCase{"foo.nadji.us", false},
		testCase{"nadji.us", false},
		testCase{"*.biz", false},
		testCase{"onizuka.homelinux.org", true},
	}
	for _, tc := range testCases {
		actual := root.ExactMatch(tc.domain)
		if tc.match != actual {
			t.Fatalf("Failed for %v (got %v expected %v)", tc.domain, actual, tc.match)
		}
	}
}

func TestWildcardMatch(t *testing.T) {
	type testCase struct {
		domain string
		match  bool
	}
	root := MakeTrie([]string{"*.google.com", "www.google.org", "*.biz", "notarealdomain", "*nadji.us", "onizuka.homelinux.org"})

	testCases := []testCase{
		testCase{"www.google.org", true},
		testCase{"www.google.com", true},
		testCase{"google.com", false},
		testCase{"google.biz", true},
		testCase{"foo.google.biz", true},
		testCase{"bar.foo.google.biz", true},
		testCase{"notarealdomain", false},
		testCase{"foo.nadji.us", false},
		testCase{"nadji.us", false},
		testCase{"*.biz", false},
		testCase{"onizuka.homelinux.org", true},
	}
	for _, tc := range testCases {
		actual := root.WildcardMatch(tc.domain)
		if tc.match != actual {
			t.Fatalf("Failed for %v (got %v expected %v)", tc.domain, actual, tc.match)
		}
	}
}
