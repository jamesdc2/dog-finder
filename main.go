package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type Dog struct {
	Name   string
	Breed  string
	Age    string
	Gender string
	Size   string
}

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("lostdogrescue.org"),
	)

	dogs := make([]Dog, 0, 200)

	// On every a element which has href attribute call callback
	c.OnHTML("ul.dog-list li", func(e *colly.HTMLElement) {

		name := e.ChildText("h3")

		dog := Dog{Name: name}

		e.ForEach("span.detail", func(_ int, el *colly.HTMLElement) {
			switch el.ChildText("h4") {
			case "Breed":
				var breeds []string
				el.DOM.Find("a").Each(func(i int, s *goquery.Selection) {
					breeds = append(breeds, s.Text())
				})
				dog.Breed = strings.Join(breeds, ", ")
			case "Size":
				dog.Size = el.DOM.First().Children().Remove().End().Text()
			case "Age":
				dog.Age = el.DOM.First().Children().Remove().End().Text()
			case "Gender":
				dog.Gender = el.DOM.First().Children().Remove().End().Text()
			}
		})

		dogs = append(dogs, dog)

	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://lostdogrescue.org/adopt/dogs-for-adoption/?age[]=young,baby")

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	enc.Encode(dogs)
}
