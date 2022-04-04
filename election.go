package blt

type Election struct {
	NumCandidates int
	NumSeats      int
	Withdrawn     []int
	Ballots       []Ballot
	Candidates    []string
	Title         string
}

type Ballot struct {
	Count       int
	Preferences [][]int
}
