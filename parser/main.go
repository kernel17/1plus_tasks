package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/playwright-community/playwright-go"
)

type data struct {
	Rank      string
	Username  string
	Title     string
	Topics    []string
	Followers string
	Country   string
	Authentic string
	Average   string
}

func main() {
	err := playwright.Install()
	var d []data
	if err != nil {
		log.Fatalln("can't install playwright's drivers")
	}

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	page.Goto("https://hypeauditor.com/top-instagram-all-russia/")
	selector, _ := page.QuerySelectorAll(".row")

	for i := 1; i < len(selector); i++ {
		var acc data

		cells, _ := selector[i].QuerySelectorAll(".row-cell")
		acc.Rank, _ = cells[0].InnerText()
		topics, _ := cells[2].QuerySelectorAll(".topic")
		if len(topics) != 0 {
			for _, topic := range topics {
				txt, _ := topic.InnerText()
				acc.Topics = append(acc.Topics, txt)
			}
		}
		contribName, _ := cells[1].QuerySelectorAll(".contributor__name-content")
		contribTitle, _ := cells[1].QuerySelectorAll(".contributor__title")
		flwrs := cells[3]
		cntr := cells[4]
		auth := cells[5]
		avg := cells[6]

		name, _ := contribName[0].InnerText()
		title, _ := contribTitle[0].InnerText()
		followers, _ := flwrs.InnerText()
		country, _ := cntr.InnerText()
		authentic, _ := auth.InnerText()
		average, _ := avg.InnerText()
		acc.Username, acc.Title, acc.Followers, acc.Country, acc.Authentic, acc.Average = name, title, followers, country, authentic, average
		d = append(d, acc)
	}
	f, err := os.Create("table.csv")
	if err != nil {
		log.Fatalf(err.Error())
	}
	for _, elem := range d {
		str := strings.Join([]string{
			elem.Rank,
			elem.Username,
			elem.Title,
			strings.Join(elem.Topics, " "),
			elem.Followers,
			elem.Country,
			elem.Authentic,
			elem.Average,
		}, ";")
		str = strings.Join([]string{str, ";"}, "")
		fmt.Fprintln(f, str)

	}

}
