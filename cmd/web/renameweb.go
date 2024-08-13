package main

import (
	"bytes"
	"html"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/phques/gorenum/renumfield"
	"github.com/phques/goweb/webhelper"
)

func homeHandler(c webhelper.Context) error {
	return c.Render("index.html", nil)
}

func renumHandler(c webhelper.Context) error {
	// read input data from the form
	text := c.FormValue("input-text")
	initialIdStr := c.FormValue("initial-id")

	// validate initial field id
	initialId, err := strconv.Atoi(initialIdStr)
	if err != nil {
		log.Printf("renumHandler, initialId must be an integer, received [%s]\n", initialIdStr)
		return webhelper.NewError(http.StatusUnprocessableEntity, "Error, initialId must be an integer")
	}

	if initialId <= 0 {
		log.Printf("renumHandler, initialId should be > 0, received [%s]\n", initialIdStr)
		return webhelper.NewError(http.StatusUnprocessableEntity, "Error, initialId should be > 0")
	}

	// create a Renum w. the text & initialId
	reader := strings.NewReader(text)
	renum := renumfield.NewRenum(initialId, reader)

	// Renumber
	log.Printf("Renumbering %d lines, starting at field id %d\n", renum.NbLines(), initialId)
	newlines := renum.Renumerate()

	// save results back to the Response
	// create one slice that holds all the lines,
	// TODO: will this be ok on Windows? (might need "\n\r")
	data := bytes.Join(newlines, []byte("\n"))
	// escape html chars
	data = []byte(html.EscapeString(string(data)))
	c.Write(data)

	return nil
}

func main() {
	// create our server
	templates := template.Must(template.ParseGlob("templates/*.html"))
	server := webhelper.CreateServer(":8080", templates)

	// add a mdw that logs what we receive
	server.AddMiddlware(func(c webhelper.Context) error {
		log.Printf("received [%s] at [%s]\n", c.Method(), c.Path())
		return nil
	})

	// register our handlers
	server.Get("/renum/home", homeHandler)
	server.Post("/renum/transform", renumHandler)

	// run the server
	log.Fatal(server.Run())
}
