package actions

import "testing"

var instance = &TitleExtract{}

var testTable = []struct {
	input  string
	expect string
}{
	{"http://google.com", "Google"},
	{"http://linux.org", "Linux.org"},
}

func TestExtract(t *testing.T) {
	for _, tt := range testTable {
		actual, err := extractTitle(tt.input)
		if err != nil {
			t.Error(err)
		}
		if actual != tt.expect {
			t.Errorf("input %s\nexpected %s\nactual %s\n", tt.input, tt.expect, actual)
		}
	}
}
