package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/boltdb/bolt"
	"github.com/is8ac/link-crawler/linkdb"
)

func getLinks(bucket *bolt.Bucket, svc *dynamodb.DynamoDB, page string) []string {
	links, err := boltReadLinks(bucket, page)
	if err != nil {
		fmt.Println(page, " is not in bolt")
		links, err = linkdb.ReadLinks(svc, page)
		if err != nil {
			fmt.Println("scraping", page)
			links = scrapeLinks(page)
			boltWriteLinks(bucket, page, links)
			time.Sleep(time.Duration(rand.Int31n(8000)) * time.Millisecond)
		}
	}
	return linkdb.StringFilter(links)
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

func boltWriteLinks(bucket *bolt.Bucket, page string, links []string) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(links)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	data := buf.Bytes()
	err = bucket.Put([]byte(page), data)
	if err != nil {
		log.Println(err)
	}
}

func boltReadLinks(bucket *bolt.Bucket, page string) (links []string, err error) {
	val := bucket.Get([]byte(page))
	if len(val) == 0 {
		err = errors.New("empty")
		return
	}
	dec := gob.NewDecoder(bytes.NewReader(val))

	err = dec.Decode(&links)
	return
}

func main() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}
	svc := dynamodb.New(sess)
	_ = svc
	db, err := bolt.Open("./bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// store some data
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("links"))
		if err != nil {
			return err
		}
		for _, page := range os.Args[1:] {
			_ = page
			recursScrape(svc, bucket, page, 3)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
