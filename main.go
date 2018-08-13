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

	// ottengo la pagina wikipedia contenente i magic number
	resp, err := http.Get("https://en.wikipedia.org/wiki/List_of_file_signatures")
	if err != nil {
		log.Fatal(err)
	}

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
	n, _ := html.Parse(strings.NewReader(`
		<th>Hex signature
</th>
<th>ISO 8859-1
</th>
<th>Offset
</th>
<th>File extension
</th>
<th>Description
</th></tr>
<tr>
<td><pre>a1 b2 c3 d4</pre>
<pre>d4 c3 b2 a1</pre>
</td>
<td><pre>....</pre>
</td>
<td>0
</td>
<td>pcap
</td>
<td>Libpcap File Format<sup id="cite_ref-1" class="reference"><a href="#cite_note-1">&#91;1&#93;</a></sup>
</td>`))
	fmt.Printf("%#v", ParseRow(n))

}

func ParseCell(n *html.Node, column int, r *Row) {

	if n.Type == html.TextNode {

	}
}

func ParseRow(n *html.Node) Row {
	var r Row
	colnum := 0
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		//fmt.Printf("%#v\n\n", c)
		colnum++
		switch {
		case colnum == 1:
			r.signature = ExtractText(c)
			fmt.Printf("Signature: -- %#v\n\n", r.signature)
		case colnum == 2:
			r.iso = ExtractText(c)
		case colnum == 3:
			r.offset = ExtractText(c)
		case colnum == 4:
			r.extension = ExtractList(c)
		case colnum == 5:
			r.description = ExtractText(c)
		}
		if colnum > 6 {
			fmt.Println("Error: too many columns")
		}
	}
	return r
}

func ExtractText(n *html.Node) []string {
	if n.Type == html.TextNode {
		trimmed := strings.TrimSpace(n.Data)
		if trimmed == "" {
			return []string{}
		}
		//trim any invisible character
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
	if n.Type == html.TextNode {
		trimmed := strings.TrimSpace(n.Data)
		if trimmed == "" {
			return []string{}
		}

		//trim some invisible character and made a list
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
