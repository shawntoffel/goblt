package blt

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	result, err := parse("simple.blt")
	if err != nil {
		t.Error(err)
	}

	assertInt(t, 6, result.NumCandidates, "candidates")
	assertInt(t, 3, result.NumSeats, "seats")
	assertIntSlice(t, []int{6}, result.Withdrawn, "withdrawn")

	assertInt(t, 5, len(result.Ballots), "ballots")

	ballots := result.Ballots
	assert2dIntSlice(t, [][]int{{1, 2}, {3}}, ballots[0].Preferences, "preferences")
	assert2dIntSlice(t, [][]int{{2}, {1}, {3}}, ballots[1].Preferences, "preferences")
	assert2dIntSlice(t, [][]int{{3}}, ballots[2].Preferences, "preferences")
	assert2dIntSlice(t, [][]int{{4}}, ballots[3].Preferences, "preferences")
	assert2dIntSlice(t, [][]int{{5}}, ballots[4].Preferences, "preferences")

	assertInt(t, 6, len(result.Candidates), "candidate names")
	assertString(t, "Test", result.Title, "title")
}

func TestNamedBallots(t *testing.T) {
	result, err := parse("simple.blt")
	if err != nil {
		t.Error(err)
	}

	named := result.NamedBallots()
	assertInt(t, 5, len(named), "ballots")
	assertInt(t, 2, len(named[0].Preferences), "preference")
	assertInt(t, 3, len(named[1].Preferences), "preference")
	assertInt(t, 1, len(named[2].Preferences), "preference")
	assertInt(t, 1, len(named[3].Preferences), "preference")
	assertInt(t, 1, len(named[4].Preferences), "preference")
}

func TestFlatNamedBallots(t *testing.T) {
	result, err := parse("simple.blt")
	if err != nil {
		t.Error(err)
	}

	named := result.FlatNamedBallots()
	assertInt(t, 5, len(named), "ballots")
	assertInt(t, 3, len(named[0].Preferences), "preference")
	assertInt(t, 3, len(named[1].Preferences), "preference")
	assertInt(t, 1, len(named[2].Preferences), "preference")
	assertInt(t, 1, len(named[3].Preferences), "preference")
	assertInt(t, 1, len(named[4].Preferences), "preference")
}

func TestNamedWithdrawn(t *testing.T) {
	result, err := parse("simple.blt")
	if err != nil {
		t.Error(err)
	}

	named := result.NamedWithdrawn()
	assertInt(t, 1, len(named), "withdrawn")
	assertString(t, "Frank", named[0], "withdrawn name")

}

func parse(filename string) (*Election, error) {
	data := loadFileData(filename)
	p := NewParser(strings.NewReader(data))
	return p.Parse()
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

func assert2dIntSlice(t *testing.T, expected, got [][]int, name string) {
	for i, val := range got {
		assertIntSlice(t, val, expected[i], name)
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
