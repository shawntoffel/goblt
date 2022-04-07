package blt

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

type Parser struct {
	buf *bufio.Reader
}

// NewParser initializes a new Parser using the provided io.Reader.
func NewParser(r io.Reader) *Parser {
	return &Parser{
		buf: bufio.NewReader(r),
	}
}

// Parse reads BLT data and returns an Election.
func (p *Parser) Parse() (*Election, error) {
	numCandidates, err := p.readInt()
	if err != nil {
		return nil, err
	}

	numSeats, err := p.readInt()
	if err != nil {
		return nil, err
	}

	withdrawn, err := p.withdrawn()
	if err != nil {
		return nil, err
	}

	ballots, err := p.ballots()
	if err != nil {
		return nil, err
	}

	candidates, err := p.candidates(numCandidates)
	if err != nil {
		return nil, err
	}

	title, err := p.title()
	if err != nil {
		return nil, err
	}

	return &Election{
		NumCandidates: numCandidates,
		NumSeats:      numSeats,
		Withdrawn:     withdrawn,
		Ballots:       ballots,
		Candidates:    candidates,
		Title:         title,
	}, nil
}

func (p *Parser) withdrawn() ([]int, error) {
	withdrawn := []int{}

	for {
		r, _, err := p.buf.ReadRune()
		if err != nil {
			return nil, err
		}

		if r == '-' {
			val, err := p.readInt()
			if err != nil {
				return nil, err
			}

			withdrawn = append(withdrawn, val)
		} else {
			err := p.buf.UnreadRune()
			if err != nil {
				return nil, err
			}

			break
		}
	}

	return withdrawn, nil
}

func (p *Parser) ballots() ([]Ballot, error) {
	ballots := []Ballot{}
	for {
		r, _, err := p.buf.ReadRune()
		if err != nil {
			return nil, err
		}

		// end of all ballot lines
		if r == '0' {
			err = p.whitespace()
			if err != nil {
				return nil, err
			}

			break
		} else {
			err := p.buf.UnreadRune()
			if err != nil {
				return nil, err
			}
		}

		ballot, err := p.ballot()
		if err != nil {
			return nil, err
		}

		ballots = append(ballots, *ballot)
	}

	return ballots, nil
}

func (p *Parser) ballot() (*Ballot, error) {
	b := &Ballot{}
	index := 0
	equal := false
	for {
		val, err := p.readInt()
		if err != nil {
			return nil, err
		}

		// end of ballot line
		if val == 0 {
			break
		}

		// start of ballot line
		if b.Count == 0 {
			b.Count = val
		} else {
			if equal && (index-1) < len(b.Preferences) {
				b.Preferences[index-1] = append(b.Preferences[index-1], val)
			} else {
				b.Preferences = append(b.Preferences, []int{val})
				index++
			}

			r, _, err := p.buf.ReadRune()
			if err != nil {
				return nil, err
			}

			if r == '=' {
				equal = true
			} else {
				equal = false
				err := p.buf.UnreadRune()
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return b, nil
}

func (p *Parser) candidates(num int) ([]string, error) {
	names := []string{}
	for i := 0; i < num; i++ {
		r, _, err := p.buf.ReadRune()
		if err != nil {
			return nil, err
		}

		if r == '"' {
			name, err := p.quotedString()
			if err != nil {
				return nil, err
			}

			names = append(names, name)
		} else {
			err := p.buf.UnreadRune()
			if err != nil {
				return nil, err
			}
		}
	}

	return names, nil
}

func (p *Parser) title() (string, error) {
	r, _, err := p.buf.ReadRune()
	if err != nil {
		return "", err
	}

	if r == '"' {
		return p.quotedString()
	} else {
		err := p.buf.UnreadRune()
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

func (p *Parser) quotedString() (string, error) {
	content := ""

	for {
		r, _, err := p.buf.ReadRune()
		if err != nil {
			return "", err
		}

		if r == '"' {
			break
		}

		content += string(r)
	}

	err := p.whitespace()
	if err != nil {
		return "", err
	}

	return content, nil
}

func (p *Parser) readInt() (int, error) {
	number := 0
	index := 0
	for {
		r, _, err := p.buf.ReadRune()
		if err != nil {
			return -1, err
		}

		if unicode.IsDigit(r) {
			number = (number * 10) + int(r-'0')
		} else {
			err := p.buf.UnreadRune()
			if err != nil {
				return -1, err
			}

			if index == 0 {
				return -1, fmt.Errorf("unexpected character '%s'", string(r))
			}
			break
		}
		index++
	}

	err := p.whitespace()
	if err != nil {
		return -1, err
	}

	return number, nil
}

func (p *Parser) whitespace() error {
	for {
		if p.eof() {
			break
		}

		r, _, err := p.buf.ReadRune()
		if err != nil {
			return err
		}

		if !unicode.IsSpace(r) {
			err := p.buf.UnreadRune()
			if err != nil {
				return err
			}
			break
		}
	}

	return nil
}

func (p *Parser) eof() bool {
	bytes, _ := p.buf.Peek(1)
	return len(bytes) == 0
}
