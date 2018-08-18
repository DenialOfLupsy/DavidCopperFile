package main

import (
	"bytes"
	"flag"
	"fmt"
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

// Flags.
var infile = flag.String("i", "", "The input file, optional if file is last parameter")

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
	//fmt.Printf("magic table : %#v", magicTable)

	// Read file.
	var currentOffset int64

	file, err := os.OpenFile("/tmp/test", os.O_RDONLY, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//var currentMagicByte string
	for _, row := range magicTable {
		for _, off := range row.offset { //[]interface{row.description, row.extension, row.iso, ro} {
			// establish what kind of value it is and transform that to an int? hex?.

			currentOffset, err = strconv.ParseInt(off, 0, 0)

			switch {
			case err == nil:
				// It's an int.
				for _, s := range row.signature {
					if strings.Contains(s, "?") {
						// do something to handle this :S

					}

					// Seek to the right offset.
					//fmt.Printf("%d\n", currentOffset)
					_, err = file.Seek(currentOffset, 0)
					if err != nil {
						log.Fatal(err)
					}

					// Elaborate the signature string.
					var s2b []byte
					temp1 := strings.ToLower(s)
					temp2 := strings.Split(temp1, " ")
					for _, v := range temp2 {
						temp3, _ := strconv.ParseInt(v, 16, 0)
						s2b = append(s2b, byte(temp3))
					}

					// Read len(s2s) byte from file.
					b := make([]byte, len(s2b))
					_, err := file.Read(b)
					if err != nil {
						log.Fatal(err)
					}

					fmt.Printf("Bytes read from file: %#v\n", b)
					fmt.Printf("Bytes in db: %#v\n", s2b)
					fmt.Printf("String in db: %#v\n", s)

					if bytes.Equal(s2b, b) {
						fmt.Printf("Filetype: %#v.\n", row.description)
						//fmt.Printf("%#v", row)
						return
					}
					//fmt.Println("no match")
				}

			case off == "any":

				// TODO make a function for what follows...... ()
				currentOffset = 0
				for currentOffset = 0; currentOffset < 100; currentOffset++ {
					for _, s := range row.signature {
						if strings.Contains(s, "?") {
							// do something to handle this :/
						}

						// Seek to the right offset.
						_, err = file.Seek(currentOffset, 0)
						if err != nil {
							log.Fatal(err)
						}

						// Elaborate the signature string.
						var s2b []byte
						temp1 := strings.ToLower(s)
						temp2 := strings.Split(temp1, " ")
						for _, v := range temp2 {
							temp3, _ := strconv.ParseInt(v, 16, 0)
							s2b = append(s2b, byte(temp3))
						}

						// Read len(s2s) byte from file.
						b := make([]byte, len(s2b))
						_, err := file.Read(b)
						if err != nil {
							log.Fatal(err)
						}

						//fmt.Printf("Bytes read from file: %#v\n", b)
						//fmt.Printf("Bytes in db: %#v\n", s2b)
						//fmt.Printf("String in db: %#v\n", s)

						if bytes.Equal(s2b, b) {
							fmt.Printf("Filetype: %#v.\n", row.description)
							//fmt.Printf("%#v", row)
							return
						}
						//fmt.Println("no match")
					}
				}
			default:
				// In any other case: TODO.
			}
		}
	}
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
	//fmt.Println(r)
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

func ElaborateOffset(s []string) {

}
