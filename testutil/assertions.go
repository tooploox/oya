package testutil

import (
	"reflect"
	"sort"
	"testing"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/diff"
)

func AssertNoErr(t *testing.T, err error, msg string, args ...interface{}) {
	if err != nil {
		t.Fatalf(errors.Wrapf(err, msg, args...).Error())
	}
}

// AssertStringsMatch compares string slices after sorting them.
func AssertStringsMatch(t *testing.T, expected []string, actual []string, msg string, args ...interface{}) {
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

func AssertObjectsEqual(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Objects are not equal. Diff:\n %v", diff.ObjectGoPrintSideBySide(expected, actual))
	}

}
