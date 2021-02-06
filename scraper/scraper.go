package main

import (
	"encoding/json"
	"fmt"
	"grass_scraper/util"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly/v2"
)

var floraDescriptors = []string {
	"Habit",
	"Inflorescence",
	"Fertile",
	"Glumes",
	"Florets",
	"Flower",
	"Fruit",
	"Distribution",
	"Notes",
}

type floraMap map[string]string

func main() {
	c := colly.NewCollector(
		/*colly.AllowedDomains("kew.org"),*/
		colly.CacheDir("./cache/kew_cache"),
	)
	detailCollector := c.Clone()

	floraMaps := make(map[uint8][]floraMap)

	// a[href]
	c.OnHTML("div.four-col.content-head-1 ul li em a[href]", func(e *colly.HTMLElement) {
		detailCollector.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("c Visiting", r.URL)
		time.Sleep(1 * time.Second)
	})

	detailCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("detailCollector visiting", r.URL)
		time.Sleep(500 * time.Millisecond)
	})

	detailCollector.OnHTML("div#BodyText", func(element *colly.HTMLElement){
		fmt.Println("detailCollector.OnHTML")
		scientificName := element.ChildText("h1")
		childTexts := element.ChildTexts("p")
		var intermediary floraMap
		intermediary = make(floraMap)

		if scientificName != "" {
			intermediary["Name"] = scientificName
		}

		for _, val := range childTexts {
			spaceIdx := util.GetCharAt(val, ' ')
			if spaceIdx == -1 {
				continue
			} else
			{
				title := val[0:spaceIdx]
				if util.Contains(floraDescriptors, title, true) {
					if title == "FERTILE" {
						title = "FERTILE SPIKELETS"
					}
					intermediary[title] = val
				}
			}
		}

		if scientificName != "" {
			startingChar := scientificName[0]
			floraMaps[startingChar] = append(floraMaps[startingChar], intermediary)
		}
	})

	c.Visit("https://www.kew.org/data/grasses-db/sppindex.htm")
	for k, v := range floraMaps {
		file, err := os.Create(fmt.Sprintf("../data/floraData_%s.json", string(k)))

		if err != nil {
			log.Fatal(err)
		}

		jsonData, _ := json.Marshal(v)
		_, err2 := file.WriteString(string(jsonData))

		if err2 != nil {
			log.Fatal(err2)
		}

		file.Close()
	}
}
///**/