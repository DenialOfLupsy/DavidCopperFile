package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Row struct {
	// TODO type
	signature, iso, offset, extension, description []string
}

var magicTable []Row

func main() {

	// Getting the wikipedia page
	resp, err := http.Get("https://en.wikipedia.org/wiki/List_of_file_signatures")
	if err != nil {
		log.Fatal(err)
	}
	n, _ := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse table with magic numbers
	ParseTable(n)

	fmt.Printf("magic table : %#v", magicTable)

	//
	//
	//
	//
	//
	//DA CANCELLARE
	// parse della pagina per ottenere la tabella
	z := html.NewTokenizer(resp.Body)
	for {
		break
		tt := z.Next()
		//fmt.Printf("%#v\n", tt)
		if tt == html.ErrorToken {
			// ...
			return
		}
		// Process the current token.
		if tt == html.StartTagToken {
			tagname, _ := z.TagName()
			//fmt.Println( string(tagname))
			_, value, _ := z.TagAttr()
			//fmt.Printf("%s", value) //!!!!!! non mi da quello che voglio, vedi output
			if string(tagname) == "table" && string(value) == "wikitable sortable" {
				//var tsignature, tiso, toffset, textension, tdescription string
				for {
					// raggiungere il tag <tr>
					tt = z.Next()
					tagname, _ = z.TagName()
					if tt == html.StartTagToken && string(tagname) == "tr" {
						//fmt.Println(z.Text())
						//fmt.Println(html.TextToken)
					}

				}
				//tt:= z.Next() //cerco il prossimo <tr>

				//magicTable = make map[string]Row
				//magicTable[] = Row{ ; ; ; ;}
			} //TODO testare la printf per vedere se stampa "wikitable sortable", poi capire come fare a prendere i 5 tr successivi e poi metterli in variabili temporanee e parsarli nel modo giusto. Nel caso in cui ci fossero più stringhe devo fare più righe della mappa?
		}
	}
	//n := html.Parse(resp.Body)
	//n, _ = html.Parse(strings.NewReader(`<html><head></head><body><table><tbody><tr><td><pre>a1 b2 c3 d4</pre><pre>d4 c3 b2 a1</pre></td><td><pre>....</pre></td><td>0</td><td>pcap</td><td>Libpcap File Format<sup id="cite_ref-1" class="reference"><a href="#cite_note-1">&#91;1&#93;</a></sup></td></tr></tbody></table></body></html>`))
	//fmt.Printf("%#v", ParseRow(n))
}

func ParseCell(n *html.Node, column int, r *Row) {
	// TODO
	if n.Type == html.TextNode {

	}
}

func ParseTable(n *html.Node) {
	/* Parse the page until the table tag is reached.
	The class tag value should be equal to "wikitable sortable", i.e. <table class="wikitable sortable"> */

	//fmt.Printf("%#v --------", n.Data)
	if n.Type == html.ElementNode && n.Data == "table" {
		for _, a := range n.Attr {
			//fmt.Println(a.Key, a.Val)
			if a.Key == "class" && a.Val == "wikitable sortable" {
				// Parse Row.
				//fmt.Println(n.Data)
				//fmt.Println(n.FirstChild.Data)
				//fmt.Println(n.FirstChild.NextSibling.Data)
				//fmt.Println(n.FirstChild.NextSibling.FirstChild.Data)
				//fmt.Printf("%#v\n", ParseRow(n.FirstChild.NextSibling.FirstChild))
				for c := n.FirstChild.NextSibling.FirstChild; c != nil; c = c.NextSibling {
					//fmt.Println(c.Data)
					magicTable = append(magicTable, ParseRow(c))
					//fmt.Println(t)
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
	var r Row
	colnum := 0

	//fmt.Println(n.Data)
	//for c := n.FirstChild.FirstChild.NextSibling.FirstChild.FirstChild.FirstChild.FirstChild; c != nil; c = c.NextSibling {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// If the type is text, then skip (it's some \n in the html). The next sibling should be a valid tag.
		if c.Type == html.TextNode {
			continue
		}
		//fmt.Println(c.Type)
		//fmt.Println(121)
		//fmt.Printf("ParseRow : %#v ..........\n\n", c.Data)
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
	//fmt.Printf("%\v", n)
	if n.Type == html.TextNode {
		trimmed := strings.TrimSpace(n.Data)
		//fmt.Printf("%\v", trimmed)
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
	// TODO: NOT TESTED YET
	if n.Type == html.TextNode {
		trimmed := strings.TrimSpace(n.Data)
		if trimmed == "" {
			return []string{}
		}

		// Removes some invisible character and makes a list.
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
