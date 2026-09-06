package service

import (
	"fmt"
	"strings"
)

type editOp int

const (
	opEqual editOp = iota
	opDelete
	opInsert
)

type diffEdit struct {
	op       editOp
	origLine int
	modLine  int
	text     string
}

// GenerateUnifiedDiff generates clean, ANSI-colorized unified diffs showing line numbers
// (@@ -line,count +line,count @@), deleted lines (-), and added lines (+).
func GenerateUnifiedDiff(filename string, original []string, modified []string) string {
	edits := computeDiffEdits(original, modified)
	if len(edits) == 0 {
		return ""
	}

	hasChanges := false
	for _, e := range edits {
		if e.op != opEqual {
			hasChanges = true
			break
		}
	}
	if !hasChanges {
		return ""
	}

	hunks := buildHunks(edits, 3)
	if len(hunks) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("--- a/%s\n", filename))
	sb.WriteString(fmt.Sprintf("+++ b/%s\n", filename))

	for _, hunk := range hunks {
		sb.WriteString(fmt.Sprintf("\033[36m@@ -%d,%d +%d,%d @@\033[0m\n",
			hunk.startOrig, hunk.countOrig, hunk.startMod, hunk.countMod))
		for _, e := range hunk.edits {
			switch e.op {
			case opEqual:
				sb.WriteString(" " + e.text + "\n")
			case opDelete:
				sb.WriteString("\033[31m-" + e.text + "\033[0m\n")
			case opInsert:
				sb.WriteString("\033[32m+" + e.text + "\033[0m\n")
			}
		}
	}

	return sb.String()
}

type hunkData struct {
	startOrig int
	countOrig int
	startMod  int
	countMod  int
	edits     []diffEdit
}

func computeDiffEdits(a, b []string) []diffEdit {
	n := len(a)
	m := len(b)
	maxD := n + m
	if maxD == 0 {
		return nil
	}

	v := make(map[int]int)
	v[1] = 0
	trace := make([]map[int]int, 0, maxD+1)

	for d := 0; d <= maxD; d++ {
		vCopy := make(map[int]int, len(v))
		for k, val := range v {
			vCopy[k] = val
		}
		trace = append(trace, vCopy)

		for k := -d; k <= d; k += 2 {
			var x int
			if k == -d || (k != d && v[k-1] < v[k+1]) {
				x = v[k+1]
			} else {
				x = v[k-1] + 1
			}
			y := x - k
			for x < n && y < m && a[x] == b[y] {
				x++
				y++
			}
			v[k] = x
			if x >= n && y >= m {
				return backtrackTrace(trace, a, b, n, m)
			}
		}
	}
	return nil
}

func backtrackTrace(trace []map[int]int, a, b []string, n, m int) []diffEdit {
	x, y := n, m
	var revEdits []diffEdit

	origLine, modLine := n, m

	for d := len(trace) - 1; d >= 0; d-- {
		vPrev := trace[d]
		k := x - y

		var prevK int
		if k == -d || (k != d && vPrev[k-1] < vPrev[k+1]) {
			prevK = k + 1
		} else {
			prevK = k - 1
		}

		prevX := vPrev[prevK]
		prevY := prevX - prevK

		for x > prevX && y > prevY {
			x--
			y--
			revEdits = append(revEdits, diffEdit{
				op:       opEqual,
				origLine: origLine,
				modLine:  modLine,
				text:     a[x],
			})
			origLine--
			modLine--
		}

		if d > 0 {
			if x == prevX {
				// Inserted b[prevY]
				revEdits = append(revEdits, diffEdit{
					op:      opInsert,
					modLine: modLine,
					text:    b[prevY],
				})
				modLine--
				y = prevY
			} else if y == prevY {
				// Deleted a[prevX]
				revEdits = append(revEdits, diffEdit{
					op:       opDelete,
					origLine: origLine,
					text:     a[prevX],
				})
				origLine--
				x = prevX
			}
		}
	}

	edits := make([]diffEdit, len(revEdits))
	for i, e := range revEdits {
		edits[len(revEdits)-1-i] = e
	}

	currOrig, currMod := 1, 1
	for i := range edits {
		switch edits[i].op {
		case opEqual:
			edits[i].origLine = currOrig
			edits[i].modLine = currMod
			currOrig++
			currMod++
		case opDelete:
			edits[i].origLine = currOrig
			currOrig++
		case opInsert:
			edits[i].modLine = currMod
			currMod++
		}
	}

	return edits
}

func buildHunks(edits []diffEdit, contextLines int) []hunkData {
	var changeIndices []int
	for i, e := range edits {
		if e.op != opEqual {
			changeIndices = append(changeIndices, i)
		}
	}
	if len(changeIndices) == 0 {
		return nil
	}

	type rangeIdx struct {
		start int
		end   int
	}

	var ranges []rangeIdx
	for _, c := range changeIndices {
		st := c - contextLines
		if st < 0 {
			st = 0
		}
		en := c + contextLines
		if en >= len(edits) {
			en = len(edits) - 1
		}
		if len(ranges) == 0 {
			ranges = append(ranges, rangeIdx{st, en})
		} else {
			last := &ranges[len(ranges)-1]
			if st <= last.end+1 {
				if en > last.end {
					last.end = en
				}
			} else {
				ranges = append(ranges, rangeIdx{st, en})
			}
		}
	}

	var hunks []hunkData
	for _, r := range ranges {
		hunkEdits := edits[r.start : r.end+1]
		countOrig, countMod := 0, 0
		startOrig, startMod := 0, 0

		for _, e := range hunkEdits {
			switch e.op {
			case opEqual:
				countOrig++
				countMod++
				if startOrig == 0 {
					startOrig = e.origLine
				}
				if startMod == 0 {
					startMod = e.modLine
				}
			case opDelete:
				countOrig++
				if startOrig == 0 {
					startOrig = e.origLine
				}
			case opInsert:
				countMod++
				if startMod == 0 {
					startMod = e.modLine
				}
			}
		}

		if startOrig == 0 {
			if r.start > 0 {
				startOrig = edits[r.start-1].origLine
			} else {
				startOrig = 0
			}
		}
		if startMod == 0 {
			if r.start > 0 {
				startMod = edits[r.start-1].modLine
			} else {
				startMod = 0
			}
		}

		hunks = append(hunks, hunkData{
			startOrig: startOrig,
			countOrig: countOrig,
			startMod:  startMod,
			countMod:  countMod,
			edits:     hunkEdits,
		})
	}

	return hunks
}
