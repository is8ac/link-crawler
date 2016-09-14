package main

import (
	"encoding/csv"
	"log"
	"os"
)

func writeMapToStdout(linkMap map[string][]string) {
	var records [][]string
	records = append(records, []string{"page", "link"})

	w := csv.NewWriter(os.Stdout)
	defer w.Flush()
	for page, links := range linkMap {
		if isGoodLink(page, filterLists) {
			for _, link := range stringFilter(links) {
				if _, exists := linkMap[link]; exists {
					records = append(records, []string{cleanLink(page), cleanLink(link)})
				}
			}
		}
	}
	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
}
