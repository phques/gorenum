package renumfield

import (
	"bytes"
	"fmt"
)

type Renum struct {
	fieldId  int
	lines    [][]byte
	Newlines [][]byte
}

func MakeRenum(initialFieldId int, lines [][]byte) *Renum {
	return &Renum{
		fieldId:  initialFieldId,
		lines:    lines,
		Newlines: make([][]byte, 0, len(lines)),
	}
}

// loops through the lines, renumerating them into r.newlines
func (r *Renum) Renumerate() {
	for _, line := range r.lines {
		newline, changed := RenumLine(line, r.fieldId)
		r.Newlines = append(r.Newlines, newline)
		if changed {
			r.fieldId++
		}
	}
}

// Adds / replaces the 'field id' part of the line.
// If the line format is not recognized, the line is kept as-is
// format: "type fieldname[: fieldId][ = value]; [// comment]"
func RenumLine(line []byte, fieldId int) (newline []byte, changed bool) {
	// left, right will hold the left and right part of the linea
	// used to reconstruct the new line
	left := line
	right := line
	if len(left) == 0 {
		return line, false
	}

	// strip right of optional comment "// comment"
	// imbedded comment is not valid  ("// // oops")
	if idx := bytes.LastIndex(left, []byte("//")); idx > 0 {
		// go warns here: 'value of right will not be used': right = origText[idx:]
		left = left[:idx]
	}

	// strip right of ";"
	if idx := bytes.LastIndex(left, []byte(";")); idx > 0 {
		right = line[idx:]
		left = left[:idx]
	} else {
		// not a recognized format
		return line, false
	}

	// strip right of (optional) "="
	spaceForEqual := ""
	if idx := bytes.Index(left, []byte("=")); idx > 0 {
		right = line[idx:]
		left = left[:idx]
		spaceForEqual = " "
	}

	// strip right of ":" (might be missing)
	if idx := bytes.Index(left, []byte(":")); idx > 0 {
		// 'right' does not change!
		left = left[:idx]
	}

	// at this point, the rightmost 'word' should be the name of the member,
	// trim whitespace
	left = bytes.TrimRight(left, "\t ")

	// any text left to process?
	if len(left) == 0 {
		// not a recognized format
		return line, false
	}

	// we can now recreate the text with an inserted ":999" after the member name
	newText := fmt.Sprintf("%s:%d%s%s", left, fieldId, spaceForEqual, right)

	return []byte(newText), true
}
