// Reads text lines from the clipboard,
// adding/replacing the fieldId (":id") of an object def from a service file.
// Results are written back to the clipboard.
// Optional parameter: "-n=initialFieldId", default = 1

package main

import (
	"bytes"
	"flag"
	"fmt"

	"github.com/phques/gorenum/renumfield"
	"golang.design/x/clipboard"
)

var initialFieldId = flag.Int("n", 0, "initial field id")

//-----------

func writeLinesToClip(lines [][]byte) {
	// create one slice that holds all the lines,
	// then write this to the clipboard
	// TODO: will this be ok on Windows? (might need "\n\r")
	data := bytes.Join(lines, []byte("\n"))
	clipboard.Write(clipboard.FmtText, data)
}

func main() {

	// check clipboard
	if err := clipboard.Init(); err != nil {
		panic(err)
	}

	// check initialFieldId parameter
	flag.Parse()
	if *initialFieldId <= 0 {
		fmt.Println("no -n value, using default initial field ID of 1")
		*initialFieldId = 1
	}

	// read data from clipboard
	println("reading lines from clipboard")

	// create a Renum w. the clipboard lines
	clipBytes := clipboard.Read(clipboard.FmtText)
	reader := bytes.NewReader(clipBytes)
	renum := renumfield.MakeRenum(*initialFieldId, reader)

	// Renumerate
	fmt.Printf("renumerating %d lines, starting at field id %d\n", renum.NbLines(), *initialFieldId)
	newlines := renum.Renumerate()

	// save results back to clipboard
	writeLinesToClip(newlines)

	// also write new lines to terminal
	fmt.Printf("%d lines copied back to clipboard:\n", len(newlines))
	for _, l := range newlines {
		fmt.Printf("%s\n", l)
	}

}
