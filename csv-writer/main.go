package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/boltdb/bolt"
	"github.com/is8ac/link-crawler/linkdb"
)

func cleanLink(link string) string {
	index := strings.Index(link, "#")
	if index > 0 {
		link = link[:index]
	}
	index = strings.Index(link, "?")
	if index > 0 {
		link = link[:index]
	}
	if link[len(link)-1:] == "/" {
		link = link[0 : len(link)-1]
	}
	if link[:5] == "https" {
		link = "http" + link[5:]
	}
	return link
}

func getLinks(svc *dynamodb.DynamoDB, page string) (links []string) {
	links, err := linkdb.ReadLinks(svc, page)
	if err != nil {
		return
	}
	return linkdb.StringFilter(links)
}

func addToRecords(records *[][]string, page string, links []string) {
	for _, line := range *records {
		if line[0] == page {
			return
		}
	}
	for _, link := range links {
		*records = append(*records, []string{cleanLink(page), cleanLink(link)})
	}
}

func memberOf(list []string, link string) bool {
	for _, i := range list {
		if link == i {
			return true
		}
	}
	return false
}

func recursScrape(svc *dynamodb.DynamoDB, bucket *bolt.Bucket, page string, recursCount int) {
	if recursCount > 0 {
		var scrapedLinks []string
		links := getLinks(bucket, svc, page)
		fmt.Println(recursCount, page, "has", len(links), "links")
		boltWriteLinks(bucket, page, links)
		for _, page := range links {
			if !memberOf(scrapedLinks, page) {
				recursScrape(svc, bucket, page, recursCount-1)
			}
		}
	}
}

func main() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := dynamodb.New(sess)
	var records [][]string
	records = append(records, []string{"page", "link"})
	var scrapedLinks []string
	for _, page := range os.Args[1:] {
		//fmt.Println("l1")
		links := getLinks(svc, page)
		addToRecords(&records, page, links)
		for _, page := range links {
			if !memberOf(scrapedLinks, page) {
				//fmt.Println("l2", page)
				links := getLinks(svc, page)
				addToRecords(&records, page, links)
				scrapedLinks = append(scrapedLinks, page)
				for _, page := range links {
					if !memberOf(scrapedLinks, page) {
						//fmt.Println("l3", page)
						links := getLinks(svc, page)
						addToRecords(&records, page, links)
						scrapedLinks = append(scrapedLinks, page)
					}
				}
			}
		}
	}

	w := csv.NewWriter(os.Stdout)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
