// Renumerate the fieldId (":id") of an object def from a service file.
// object Toto
//
//	{
//		string Val1;
//		int Val2: 200;
//		byte Val3 = 19	// missing ';' !
//		string Val4 ;
//		string Val5="toto" ;
//	}
//
//	==>
//
// object Toto
//
//	{
//		string Val1:100;
//		int Val2:101;
//		byte Val3 = 19  // not changed,  missing ';'
//		string Val4:102;
//		string Val5:103 ="toto" ;
//	}
package renumfield

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
)

type Renum struct {
	fieldId int
	lines   [][]byte
}

func NewRenum(initialFieldId int, reader io.Reader) *Renum {
	// read lines from the reader
	lines := readLines(reader)

	// return a new Renum
	return &Renum{
		fieldId: initialFieldId,
		lines:   lines,
	}
}

func readLines(reader io.Reader) [][]byte {
	lines := make([][]byte, 0)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		// get the line as a string, it will do an allocation for us
		line := scanner.Text()
		lines = append(lines, []byte(line))
	}
	return lines
}

func (r *Renum) NbLines() int {
	return len(r.lines)
}

// loops through the lines, renumerating them
func (r *Renum) Renumerate() [][]byte {
	newlines := make([][]byte, 0, len(r.lines))

	for _, line := range r.lines {
		newline, changed := RenumLine(line, r.fieldId)
		newlines = append(newlines, newline)
		if changed {
			r.fieldId++
		}
	}

	return newlines
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
	rexp := regexp.MustCompile(`\s*=\s*`)
	if loc := rexp.FindIndex(left); loc != nil {
		idx := loc[0]
		right = line[idx:]
		left = left[:idx]
	}

	// strip right of ":" (might be missing)
	if idx := bytes.Index(left, []byte(":")); idx > 0 {
		// 'right' does not change! i.e. we do not include the old field id
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
