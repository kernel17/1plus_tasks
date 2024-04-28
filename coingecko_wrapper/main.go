package main

import (
	"log"
	"wrapper/client"
)

func main() {
	c := client.NewClient()
	log.Println(c.GetFromStorage("xmr"))
	log.Println(c.GetFromStorage("ltc"))
	log.Println(c.GetFromStorage("ton"))
	log.Println(c.GetFromStorage("all"))
	c.StartServer() // starts api server at :7878

}
