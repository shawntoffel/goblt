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

type NamedBallot struct {
	Count       int
	Preferences [][]string
}

type FlatNamedBallot struct {
	Count       int
	Preferences []string
}

func (e *Election) NamedWithdrawn() []string {
	namedWithdrawn := make([]string, len(e.Withdrawn))

	for i, index := range e.Withdrawn {
		namedWithdrawn[i] = e.Candidates[index-1]
	}

	return namedWithdrawn
}

func (e *Election) NamedBallots() []NamedBallot {
	namedBallots := make([]NamedBallot, len(e.Ballots))

	for i, ballot := range e.Ballots {
		namedBallot := NamedBallot{
			Count:       ballot.Count,
			Preferences: make([][]string, len(ballot.Preferences)),
		}

		for j, preferences := range ballot.Preferences {
			namedBallot.Preferences[j] = make([]string, len(preferences))

			for k, preference := range preferences {
				namedBallot.Preferences[j][k] = e.Candidates[preference-1]
			}
		}
		namedBallots[i] = namedBallot
	}

	return namedBallots
}

func (e *Election) FlatNamedBallots() []FlatNamedBallot {
	namedballots := e.NamedBallots()
	flatNamedBallots := make([]FlatNamedBallot, len(namedballots))

	for i, ballot := range namedballots {
		names := []string{}
		for _, preferences := range ballot.Preferences {
			for _, preference := range preferences {
				names = append(names, preference)
			}
		}

		flatNamedBallots[i] = FlatNamedBallot{
			Count:       ballot.Count,
			Preferences: names,
		}
	}

	return flatNamedBallots
}
