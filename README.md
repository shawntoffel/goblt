# goblt
A BLT file format parser in Go. Used for election data.

## Background

```
6 3
-2 -6
28 1=2 3 0
26 2 1 3 0
3 3 0
2 4 0
1 5 0
0
"Alice"
"Bob"
"Chris"
"Don"
"Eric"
"Frank"
"My Election"
```
* The first line `6 3` indicates the number of candidates `6` and seats `3`. 
* The second line `-2 -6` is optional and may be omitted. It contains negative candidate numbers indicating candidates 2 and 6 are withdrawn.
* The next lines contain vote data. `28 1=2 3 0` indicates there are 28 votes each containing candidate preferences `1=2 3`. The end of a vote line is indicated by a `0`. The preference `1=2 3` indicates `1` and `2` are tied for first preference, and `3` is second preference. 

* A line with a single `0` indicates the end of all vote data.
* Next are quoted candidate names listed in a specific order matching the numbers for vote preferences. 
* The last line is the title of the election.

## Use

`NewParser` accepts an `io.Reader` and initializes a new parser. 

```go
parser := blt.NewParser(reader)

election, err := parser.Parse()
```

The resulting parsed election struct will look like this:
```json
{
    "NumCandidates":6,
    "NumSeats":3,
    "Withdrawn":[2,6],
    "Ballots":[
        {
            "Count":28,
            "Preferences":[[1,2],[3]]
        },
        {
            "Count":26,
            "Preferences":[[2],[1],[3]]
        },
        {
            "Count":3,
            "Preferences":[[3]]
        },
        {
            "Count":2,
            "Preferences":[[4]]
        },
        {
            "Count":1,
            "Preferences":[[5]]
        }
    ],
    "Candidates": [
        "Alice",
        "Bob",
        "Chris",
        "Don",
        "Eric",
        "Frank"
    ],
    "Title":"Test"
}
```
