package prose

import (
	"fmt"
	"testing"
)

func TestTwoElements(t *testing.T) {
	list := []string{"a", "b"}
	want := "a and b"
	got := JoinWithCommas(list)
	if got != want {
		t.Error(errorString(list, got, want))
	}
}

func TestThreeElements(t *testing.T) {
	list := []string{"a", "b", "c"}
	want := "a, b, and c"
	got := JoinWithCommas(list)
	if got != want {
		t.Error(errorString(list, got, want))
	}
}

func errorString(list []string, got string, want string) string {
	return fmt.Sprintf("JoinWithCommas(%#v) = \"%s\",want \"%s\"", list, got, want)
}

func TestOneElement(t *testing.T) {
	list := []string{"a"}
	want := "a"
	got := JoinWithCommas(list)
	if got != want {
		t.Error(errorString(list, got, want))
	}
}
func TestZeroElement(t *testing.T) {
	list := []string{""}
	want := ""
	got := JoinWithCommas(list)
	if got != want {
		t.Error(errorString(list, got, want))
	}
}
