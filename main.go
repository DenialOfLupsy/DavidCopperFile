package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type Row struct {
	// TODO later differentiate types
	signature, iso, offset, extension, description []string
}

var magicTable []Row

func main() {
	flag.Parse()

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

	if len(flag.Args()) == 0 {
		fmt.Println("Please provide at least one file.")
	}
	for _, fname := range flag.Args() {
		file, err := os.OpenFile(fname, os.O_RDONLY, 0444)
		if err != nil {
			log.Fatal(err)
		}
		doTheMagic(file)
		file.Close()
	}

}

func doTheMagic(file io.ReaderAt) {

loop:
	for _, row := range magicTable {
		for _, off := range row.offset {

			currentOffset, err := strconv.ParseInt(off, 0, 0)

			switch {
			case err == nil:
				// It's an int.
				if Match(currentOffset, file, row.signature) {
					fmt.Printf("Filetype: %#v.\n", row.description)
					continue loop
				}

			case off == "any":
				for currentOffset = 0; currentOffset < 256; currentOffset++ {
					if Match(currentOffset, file, row.signature) {
						fmt.Printf("Filetype: %#v.\n", row.description)
						continue loop
					}
				}

			default:
				// Use this case to elaborate other offset type.
			}
		}
	}
}

// ParseTable parses the page until the table tag is reached. The class tag
// value should be equal to "wikitable sortable", i.e. <table class="wikitable sortable">
func ParseTable(n *html.Node) {

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

// ParseRow .
func ParseRow(n *html.Node) Row {
	var r Row
	colnum := 0

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// If the type is text, then skip (it's some \n in the html).
		// The next sibling should be a valid tag.
		if c.Type == html.TextNode {
			continue
		}

		colnum++

		switch colnum {
		case 1:
			r.signature = ExtractText(c)
		case 2:
			r.iso = ExtractText(c)
		case 3:
			r.offset = ExtractText(c)
		case 4:
			r.extension = ExtractText(c)
		case 5:
			r.description = ExtractText(c)
		default:
			fmt.Println("Error: too many columns.")
		}
	}
	return r
}

// ExtractText extracts text recursively from the n node down to the bottom of the tree.
func ExtractText(n *html.Node) []string {

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

// Match tests if the signatures match the ones in the current file.
func Match(currentOffset int64, file io.ReaderAt, signatures []string) bool {

	for _, s := range signatures {
		if strings.Contains(s, "?") {
			// do something to handle this :S

		}

		// Elaborate the signature string.
		var s2b []byte
		temp1 := strings.ToLower(s)
		temp2 := strings.Split(temp1, " ")
		for _, v := range temp2 {
			temp3, _ := strconv.ParseInt(v, 16, 0)
			s2b = append(s2b, byte(temp3))
		}

		b := make([]byte, len(s2b))
		_, err := file.ReadAt(b, currentOffset)
		if err == io.EOF {
			return false
		}
		if err != nil {
			log.Fatal(err)
		}

		if bytes.Equal(s2b, b) {
			return true
		}
	}
	return false
}
