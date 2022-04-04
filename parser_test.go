package blt

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	data := loadFileData("simple.blt")

	p := NewParser(strings.NewReader(data))

	result, err := p.Parse()
	if err != nil {
		t.Error(err)
	}

	assertInt(t, 6, result.NumCandidates, "candidates")
	assertInt(t, 3, result.NumSeats, "seats")
	assertIntSlice(t, []int{6}, result.Withdrawn, "withdrawn")
	assertInt(t, 5, len(result.Ballots), "ballots")
	assertInt(t, 6, len(result.Candidates), "candidate names")
	assertString(t, "Test", result.Title, "title")
}

func loadFileData(filename string) string {
	bytes, _ := ioutil.ReadFile("testdata/" + filename)
	return string(bytes)
}

func assertInt(t *testing.T, expected, got int, name string) {
	if expected != got {
		t.Errorf("expected %d %s, got %d", expected, name, got)
	}
}

func assertString(t *testing.T, expected, got string, name string) {
	if expected != got {
		t.Errorf("expected %s %s, got %s", expected, name, got)
	}
}

func assertIntSlice(t *testing.T, expected, got []int, name string) {
	if !testEqual(expected, got) {
		t.Errorf("expected %d %s, got %d", expected, name, got)
	}
}

func testEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
