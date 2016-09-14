package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

func getLinks(db *bolt.DB, page string) []string {
	links, err := boltReadLinks(db, page)
	if err != nil {
		fmt.Println("scraping", page)
		links = scrapeLinks(page)
		boltWriteLinks(db, page, links)
		time.Sleep(time.Duration(rand.Int31n(3000)) * time.Millisecond) // Be kind
	}
	return stringFilter(links)
}

func memberOf(list []string, link string) bool {
	for _, i := range list {
		if link == i {
			return true
		}
	}
	return false
}

func recursScrape(db *bolt.DB, linkMap *map[string][]string, page string, recursCount int) {
	if recursCount > 0 {
		links := getLinks(db, page)
		//fmt.Println(recursCount, page, "has", len(links), "links")
		(*linkMap)[page] = links
		for _, page := range links {
			if _, exists := (*linkMap)[page]; !exists {
				recursScrape(db, linkMap, page, recursCount-1)
			}
		}
	}
}

var filterLists = []filterList{
	{White: true, Globs: []string{
	//"*",
	}},
	{White: false, Globs: []string{
		"*archive.org*",
		"*.pdf",
		"*.csv",
		"*.txt",
		"*.png",
		"*.xls",
		"*.xlsx",
		"*.jpeg",
		"*.jpg",
		"*.svg",
		"*.ogg",
		"*wikipedia.org/wiki/Wikipedia:*",
		"*wikipedia.org/wiki/Category:*",
		"*wikipedia.org/wiki/Portal:*",
		"*wikipedia.org/wiki/Special:*",
		"*wikipedia.org/wiki/Help:*",
		"*wikipedia.org/wiki/Talk:*",
		"*wikipedia.org/wiki/File:*",
		"*wikipedia.org/wiki/Template_talk:*",
		"*wikipedia.org/wiki/Template:*",
		"*wikipedia.org/wiki/Main_Page*",
		"*slatestarcodex.com/20[0-9][0-9]/[0-9][0-9]",
		"*slatestarcodex.com/20[0-9][0-9]/",
		"*slatestarcodex.com/*terrorists-vs-chairs*",
		"*slatestarcodex.com/*open-thread-57-75",
		"*slatestarcodex.com/*open-thread-57-25",
		"*slatestarcodex.com/*open-thread-57-5",
		"*slatestarcodex.com/*reverse-*-brand-name-drugs",
		"http://en.wikipedia.org/wiki/International_Standard_Book_Number",
		"http://en.wikipedia.org/wiki/*(disambiguation)",
		//"*slatstarcodex.com//*",
		//"*analog*",
	}},
	{White: true, Globs: []string{
		//"*",
		"*en.wikipedia.org/wiki/*",
		//"*//slatestarcodex.com",
		//"*slatestarcodex.com/20[0-9][0-9]/[0-9][0-9]/*",
		//"*en.wikipedia*",
		//"*http://isaacleonard.com*",
	}},
	{White: false, Globs: []string{
		"*.wikipedia.org*",
		"*wikimedia.org*",
		"*wikimediafoundation.org*",
		"*mediawiki.org*",
		"*wikidata.org*",
		"*dbpedia.org*",
		"*amazonaws.com*",
		"*",
	}},
	{White: true, Globs: []string{"*"}},
}

func main() {
	db, err := bolt.Open("bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	linkMap := make(map[string][]string)

	// store some data
	err = db.Update(func(tx *bolt.Tx) error {
		_, createErr := tx.CreateBucketIfNotExists([]byte("links"))
		if err != nil {
			return createErr
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, page := range os.Args[1:] {
		_ = page
		recursScrape(db, &linkMap, page, 3)
	}
	//fmt.Println(len(linkMap))
	writeMapToStdout(linkMap)
}
