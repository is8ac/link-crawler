package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"

	"github.com/boltdb/bolt"
	"github.com/gobwas/glob"
)

type filterList struct {
	White bool
	Globs []string
}

func isGoodLink(link string, filterLists []filterList) (isGood bool) {
	for _, list := range filterLists {
		for _, pat := range list.Globs {
			g := glob.MustCompile(pat)
			if g.Match(link) {
				return list.White
			}
		}
	}
	return true
}

func stringFilter(strings []string) (output []string) {
	for _, item := range strings {
		if isGoodLink(item, filterLists) {
			output = append(output, item)
		}
	}
	return
}

func boltWriteLinks(db *bolt.DB, page string, links []string) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(links)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	data := buf.Bytes()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("links"))
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(page), data)
		if err != nil {
			log.Println(err)
		}
		return nil
	})

}

func boltReadLinks(db *bolt.DB, page string) (links []string, err error) {
	var val []byte
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("links"))
		val = bucket.Get([]byte(page))
		return nil
	})
	if len(val) == 0 {
		err = errors.New("empty")
		return
	}
	dec := gob.NewDecoder(bytes.NewReader(val))

	err = dec.Decode(&links)

	return
}
