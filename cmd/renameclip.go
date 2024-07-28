// Reads text lines from the clipboard,
// adding/replacing the fieldId (":id") of an object def from a service file.
// Results are written back to the clipboard.
// Optional parameter: "-n=initialFieldId", default = 1

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

func getClipAsLines() [][]byte {
	clipBytes := clipboard.Read(clipboard.FmtText)
	clip := bytes.ReplaceAll(clipBytes, []byte("\r\n"), []byte("\n"))
	lines := bytes.Split(clip, []byte("\n"))
	return lines
}

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
		fmt.Printf("no -n value, using default initial field ID of 1")
		*initialFieldId = 1
	}

	// read data from clipboard
	println("reading lines from clipboard")
	lines := getClipAsLines()

	// create a Renum w. the clipboard lines
	renum := renumfield.MakeRenum(*initialFieldId, lines)

	// Renumerate
	fmt.Printf("renumerating %d lines, starting at field id %d\n", len(lines), *initialFieldId)
	renum.Renumerate()

	// save results back to clipboard
	writeLinesToClip(renum.Newlines)

	// write new lines to terminal
	fmt.Printf("%d lines copied back to clipboard:\n", len(lines))
	for _, l := range renum.Newlines {
		fmt.Printf("%s\n", l)
	}

}
