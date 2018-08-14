package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Row struct {
	// TODO later differentiate types
	signature, iso, offset, extension, description []string
}

var magicTable []Row

func main() {

	// Getting the wikipedia page.
	resp, err := http.Get("https://en.wikipedia.org/wiki/List_of_file_signatures")
	if err != nil {
		log.Fatal(err)
	}

	n, _ := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse html and set the variable magicTable.
	ParseTable(n)
	fmt.Printf("magic table : %#v", magicTable)
}

func ParseCell(n *html.Node, column int, r *Row) {
	// TODO
	if n.Type == html.TextNode {

	}
}

func ParseTable(n *html.Node) {
	/* Parse the page until the table tag is reached.
	The class tag value should be equal to "wikitable sortable", i.e. <table class="wikitable sortable"> */

	if n.Type == html.ElementNode && n.Data == "table" {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "wikitable sortable" {
				for c := n.FirstChild.NextSibling.FirstChild; c != nil; c = c.NextSibling {
					thisRow := ParseRow(c)
					if len(thisRow.signature) == 0 &&
						len(thisRow.iso) == 0 &&
						len(thisRow.offset) == 0 &&
						len(thisRow.extension) == 0 &&
						len(thisRow.description) == 0 {
						continue
					}
					magicTable = append(magicTable, thisRow)
				}
			}
		}
	}
	// Keep looking for the right table.
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ParseTable(c)
	}
}

func ParseRow(n *html.Node) Row {
	/* Parses row differently for each column */
	var r Row
	colnum := 0

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// If the type is text, then skip (it's some \n in the html). The next sibling should be a valid tag.
		if c.Type == html.TextNode {
			continue
		}

		colnum++
		switch {
		case colnum == 1:
			r.signature = ExtractText(c)
			//fmt.Printf("Signature: -- %#v\n\n", r.signature)
		case colnum == 2:
			r.iso = ExtractText(c)
			//fmt.Printf("ISO: -- %#v\n\n", r.iso)
		case colnum == 3:
			r.offset = ExtractText(c)
			//fmt.Printf("Offset: -- %#v\n\n", r.offset)
		case colnum == 4:
			r.extension = ExtractList(c)
			//fmt.Printf("Extension: -- %#v\n\n", r.extension)
		case colnum == 5:
			r.description = ExtractText(c)
			//fmt.Printf("Description: -- %#v\n\n", r.description)
		}
		if colnum > 6 {
			fmt.Println("Error: too many columns")
		}
	}
	fmt.Println(r)
	return r
}

func ExtractText(n *html.Node) []string {
	/* Extract text recursively from the n node down to the bottom of the tree. */
	if n.Type == html.TextNode {
		trimmed := strings.TrimSpace(n.Data)
		if trimmed == "" {
			return []string{}
		}
		// Removes any invisible character.
		for _, i := range []string{"\a", "\b", "\f", "\n", "\r", "\t", "\v"} {
			trimmed = strings.Replace(trimmed, i, "", -1)
		}
		return []string{trimmed}
	}
	var accumulator []string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		accumulator = append(accumulator, ExtractText(c)...)
	}

	return accumulator
}

func ExtractList(n *html.Node) []string {
	// TODO: NOT TESTED NOR USED YET
	if n.Type == html.TextNode {
		trimmed := strings.TrimSpace(n.Data)
		if trimmed == "" {
			return []string{}
		}

		// Removes some invisible character and makes a comma separated list.
		for _, i := range []string{"\a", "\b", "\f", "\r", "\t", "\v"} {
			trimmed = strings.Replace(trimmed, i, "", -1)
		}
		trimmed = strings.Replace(trimmed, "\n", ", ", -1)
		return []string{trimmed}
	}
	var accumulator []string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		accumulator = append(accumulator, ExtractText(c)...)
	}
	return accumulator
}
