package compare

import (
	"fmt"
	"testing"
)

func TestLeftLarger(t *testing.T) {
	a := 3
	b := 2
	got := Larger(a, b)
	want := 3
	if got != want {
		t.Error(errorString(a, b, got, want))
	}
}
func TestRightLarger(t *testing.T) {
	a := 3
	b := 4
	got := Larger(a, b)
	want := 4
	if got != want {
		t.Error(errorString(a, b, got, want))
	}
}
func errorString(a int, b int, got int, want int) string {
	return fmt.Sprintf("Larger(%#v,%#v) = %d, want %d", a, b, got, want)
}
