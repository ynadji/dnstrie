package dnstrie

import (
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
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
	}

	for _, tc := range testCases {
		reversedLabels := reverseLabelSlice(tc.domain)
		if !reflect.DeepEqual(reversedLabels, tc.reversedLabels) {
			t.Fatalf("Failed to reverse labels. Got %+v expected %+v.", reversedLabels, tc.reversedLabels)
		}
	}
}
