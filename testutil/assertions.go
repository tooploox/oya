package testutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"sort"
	"testing"

	"github.com/go-test/deep"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/diff"
)

func AssertNoErr(t *testing.T, err error, msg string, args ...interface{}) {
	t.Helper()
	if err != nil {
		t.Fatalf(errors.Wrapf(err, msg, args...).Error())
	}
}

func AssertErr(t *testing.T, err error, msg string, args ...interface{}) {
	t.Helper()
	if err == nil {
		t.Fatalf(msg, args...)
	}
}

func AssertTrue(t *testing.T, b bool, msg string, args ...interface{}) {
	t.Helper()
	if !b {
		t.Fatalf(msg, args...)
	}
}

func AssertFalse(t *testing.T, b bool, msg string, args ...interface{}) {
	t.Helper()
	if b {
		t.Fatalf(msg, args...)
	}
}

// AssertStringsMatch compares string slices after sorting them.
func AssertStringsMatch(t *testing.T, expected []string, actual []string, msg string, args ...interface{}) {
	t.Helper()
	expSorted := make([]string, len(expected))
	copy(expSorted, expected)
	sort.Strings(expSorted)
	actSorted := make([]string, len(actual))
	copy(actSorted, actual)
	sort.Strings(actSorted)

	if !reflect.DeepEqual(expSorted, actSorted) {
		t.Errorf(msg, args...)
	}
}

// AssertRegexpMatch checks if string matches the regexp.
func AssertRegexpMatch(t *testing.T, expectedRegexp, actual string) {
	rx := regexp.MustCompile(expectedRegexp)
	if !rx.MatchString(actual) {
		t.Errorf("Expected regexp %q to match %q", expectedRegexp, actual)
	}
}

func AssertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %v actual: %v", expected, actual)
	}
}

func AssertObjectsEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if df := deep.Equal(expected, actual); df != nil {
		t.Errorf("Objects are not equal.\n\nDiff: expected\tactual\n %v\n\nSide-by-side: %v", df, diff.ObjectGoPrintSideBySide(expected, actual))
	}
}

func AssertObjectsEqualMsg(t *testing.T, expected, actual interface{}, msg string, args ...interface{}) {
	t.Helper()
	if df := deep.Equal(expected, actual); df != nil {
		t.Errorf("%v: %v",
			fmt.Sprintf(msg, args...),
			fmt.Sprintf("objects are not equal.\n\nDiff:\n %v\n\nSide-by-side:\n%v", df, diff.ObjectGoPrintSideBySide(expected, actual)))
	}
}

func AssertPathExists(t *testing.T, path string) {
	t.Helper()
	_, err := os.Stat(path)
	if err != nil {
		t.Errorf("path %v does not exist", path)
	}
}

func AssertFileContains(t *testing.T, path string, expectedContent string) {
	t.Helper()
	AssertPathExists(t, path)
	actual, err := ioutil.ReadFile(path)
	AssertNoErr(t, err, "Expected no error reading %v", path)
	AssertEqual(t, expectedContent, string(actual))
}
