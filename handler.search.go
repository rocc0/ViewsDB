package main

import (
	"log"

	"github.com/blevesearch/bleve"

	"github.com/gin-gonic/gin"
)


const (
	viewsIdx   = "views.index"
)

var bleveIdx bleve.Index

// Bleve connect or create the index persistence
func Bleve(indexPath string) (bleve.Index, error) {

	log.Printf("Indexing...")

	if bleveIdx == nil {
		var err error

		bleveIdx, err = bleve.Open(indexPath)

		if err != nil {

			mapping := bleve.NewIndexMapping()
			bleveIdx, err = bleve.New(indexPath, mapping)
			if err != nil {
				return nil, err
			}
		}
	}
	// return the index
	return bleveIdx, nil
}


func generateIndexes() {
	eventList := getAllArticles()
	idx := idxCreate()
	indexEvents(idx, eventList)
}


func indexEvents(idx bleve.Index, eventList []View_mainpage) {
	for _, event := range eventList {
		event.Index(idx)
	}
}

func (e *View_mainpage) Index(index bleve.Index) error {
	err := index.Index(string(e.Id), e)
	return err
}

func idxCreate() bleve.Index{
	idx, err := Bleve(viewsIdx)
	if err != nil || idx == nil {
		log.Print(err.Error())
	}
	return idx

}


func searchIndex(c *gin.Context) {
	idx := idxCreate()
	query := bleve.NewQueryStringQuery("рішенням Постанова 00:00:00 січень")
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Highlight = bleve.NewHighlight()
	searchResult, err := idx.Search(searchRequest)
	if err != nil {
		log.Print(err.Error())
	}
	for k, v := range searchResult.Hits[4].Fragments {
		log.Print("Key: ", k, " Value: ", v[0])
	}
	log.Print()
}

