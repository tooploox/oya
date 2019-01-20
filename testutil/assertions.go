package testutil

import (
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"testing"

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

func AssertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %v actual: %v", expected, actual)
	}
}

func AssertObjectsEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Objects are not equal. Diff:\n %v", diff.ObjectGoPrintSideBySide(expected, actual))
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
