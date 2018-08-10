package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

type Row struct {
	// TODO type
	signature, iso, offset, extension, description string
}

var magicTable map[string]Row

func main() {

	// ottengo la pagina wikipedia contenente i magic number
	resp, err := http.Get("https://en.wikipedia.org/wiki/List_of_file_signatures")
	if err != nil {
		log.Fatal(err)
	}

	// parse della pagina per ottenere la tabella
	z := html.NewTokenizer(resp.Body)
	for {
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
			fmt.Printf("%s", value) //!!!!!! non mi da quello che voglio, vedi output
			if string(tagname) == "table" && string(value) == "wikitable sortable" {
				//var tsignature, tiso, toffset, textension, tdescription string
				for {
					// raggiungere il tag <tr>
					tt = z.Next()
					tagname, _ = z.TagName()
					if tt == html.StartTagToken && string(tagname) == "tr" {
						fmt.Println(z.Text())
						fmt.Println(html.TextToken)
					}

				}
				//tt:= z.Next() //cerco il prossimo <tr>

				//magicTable = make map[string]Row
				//magicTable[] = Row{ ; ; ; ;}
			} //TODO testare la printf per vedere se stampa "wikitable sortable", poi capire come fare a prendere i 5 tr successivi e poi metterli in variabili temporanee e parsarli nel modo giusto. Nel caso in cui ci fossero più stringhe devo fare più righe della mappa?
		}
	}

}
