package flags_test

import (
	"testing"

	"github.com/bilus/oya/pkg/flags"
	tu "github.com/bilus/oya/testutil"
)

func TestEmpty(t *testing.T) {
	positionalArgs, flags, err := flags.Parse(nil)
	tu.AssertNoErr(t, err, "flag.Parse failed")
	if len(positionalArgs) > 0 {
		t.Errorf("Expected no positional arguments, actual: %v", len(positionalArgs))
	}
	if len(flags) > 0 {
		t.Errorf("Expected no positional arguments, actual: %v", len(flags))
	}
}

func TestPositionalArgs(t *testing.T) {
	positionalArgs, flags, err := flags.Parse([]string{"arg1", "arg2"})
	tu.AssertNoErr(t, err, "flag.Parse failed")
	if len(positionalArgs) != 2 {
		t.Errorf("Expected 2 positional arguments, actual: %v", len(positionalArgs))
	}
	tu.AssertObjectsEqual(t, []string{"arg1", "arg2"}, positionalArgs)

	if len(flags) > 0 {
		t.Errorf("Expected no positional arguments, actual: %v", len(flags))
	}
}

func TestSwitchFlags(t *testing.T) {
	positionalArgs, flags, err := flags.Parse([]string{"--switch1", "--switch2"})
	tu.AssertNoErr(t, err, "flag.Parse failed")
	if len(positionalArgs) > 0 {
		t.Errorf("Expected no positional arguments, actual: %v", len(positionalArgs))
	}
	if len(flags) != 2 {
		t.Errorf("Expected 2 flags, actual: %v", len(flags))
	}
	assertFlagSet(t, flags, "switch1")
	assertFlagSet(t, flags, "switch2")
}

// Repeating the same flag is not supported yet. This may change.
func TestRepeatedFlag(t *testing.T) {
	_, _, err := flags.Parse([]string{"--switch", "--switch"})
	if err == nil {
		t.Errorf("Expected an error for a repeated flag")
	}
}

func TestValuesFlags(t *testing.T) {
	positionalArgs, flags, err := flags.Parse([]string{"--value1=123", "--value2=42"})
	tu.AssertNoErr(t, err, "flag.Parse failed")
	if len(positionalArgs) > 0 {
		t.Errorf("Expected no positional arguments, actual: %v", len(positionalArgs))
	}
	if len(flags) != 2 {
		t.Errorf("Expected 2 flags, actual: %v", len(flags))
	}
	assertFlagEqual(t, flags, "value1", "123")
	assertFlagEqual(t, flags, "value2", "42")
}

func TestFullMix(t *testing.T) {
	positionalArgs, flags, err := flags.Parse([]string{"some-arg", "--some-value=123", "--some-switch"})
	tu.AssertNoErr(t, err, "flag.Parse failed")
	if len(positionalArgs) != 1 {
		t.Errorf("Expected 1 positional argument, actual: %v", len(positionalArgs))
	}
	if len(flags) != 2 {
		t.Errorf("Expected 2 flags, actual: %v", len(flags))
	}
	assertFlagSet(t, flags, "some-value")
	assertFlagSet(t, flags, "some-switch")
	assertFlagEqual(t, flags, "some-value", "123")
}

func assertFlagSet(t *testing.T, flags map[string]string, flag string) {
	_, ok := flags[flag]
	if !ok {
		t.Errorf("Expected flag %q to be set", flag)
	}
}

func assertFlagEqual(t *testing.T, flags map[string]string, flag, expectedValue string) {
	actualValue, ok := flags[flag]
	if !ok {
		t.Errorf("Expected flag %q to be set", flag)
	}
	if actualValue != expectedValue {
		t.Errorf("Expected flag %q to be set to %q, actual: %q", flag, expectedValue, actualValue)
	}
}
