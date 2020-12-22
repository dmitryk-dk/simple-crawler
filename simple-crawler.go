package main

import (
	"log"

	"github.com/dmitryk-dk/simlpe-crawler/crawler"
)

func main() {
	c := crawler.New("https://www.brightlocal.com", 10)
	go c.Process()
	log.Printf("visited links -> %+v", <-c.VisitedLinks())
}
