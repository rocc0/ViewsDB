package main

import (
	"log"

	"github.com/blevesearch/bleve"
	"strconv"
)

const (
	viewsIdx = "views.index"
)

var bleveIdx bleve.Index

func Bleve(indexPath string) (bleve.Index, error) {
	if bleveIdx == nil {
		var err error
		bleveIdx, err = bleve.Open(indexPath)
		if err != nil {
			mapping, _ := buildIndexMapping()
			bleveIdx, err = bleve.New(indexPath, mapping)
			if err != nil {
				log.Print(err)
			}
		}
	}
	return bleveIdx, nil
}

func generateIndexes() {
	viewList := getAllArticles()
	idx, err := Bleve(viewsIdx)
	if err != nil || idx == nil {
		log.Print(err.Error())
	}
	log.Printf("Indexing...")
	for _, view := range viewList {
		view.Index(idx)
	}

}

func (e *View_mainpage) Index(index bleve.Index) error {
	err := index.Index(string(strconv.Itoa(e.Id)), e)
	return err
}

