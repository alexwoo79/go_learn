package prose

import (
	"testing"
)

type testData struct {
	list []string
	want string
}

func TestJoinElements(t *testing.T) {
	tests := []testData{
		{list: []string{}, want: ""},
		{list: []string{"a"}, want: "a"},
		{list: []string{"a", "b"}, want: "a and b"},
		{list: []string{"a", "b", "c"}, want: "a, b, and c"},
	}
	for _, test := range tests {

		got := JoinWithCommas(test.list)
		if got != test.want {
			t.Errorf("JoinWithCommas(%#v) = \"%s\",want \"%s\"", test.list, got, test.want)
		}
	}
}
