package main

import (
	"context"
	"log"

	elastic "gopkg.in/olivere/elastic.v5"
)

type Idx struct {
	Id         string `json:"id"`
	Requisits  string `json:"requisits"`
	Reg_Date   string `json:"reg_date"`
	Gov_choice string `json:"gov_choice"`
	Trace_year string `json:"year"`
	Developer  string `json:"developer"`
	Base       string `json:"base"`
	Repeated   string `json:"repeat"`
	Periodical string `json:"period"`
	Fact       string `json:"fact"`
}


const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"trace":{
			"properties":{
				"id":{
					"type":"text"
				},
				"requisits":{
					"type":"text",
					"analyzer": "ukrainian"
				},
				"reg_date":{
					"type":"text"
				},
				"gov_choice":{
					"type":"string",
					"analyzer": "ukrainian"
				},
				"developer":{
					"type":"integer"
				},
				"year":{
					"type":"integer"
				},
				"base":{
					"type":"integer"
				},
				"repeat":{
					"type":"integer"
				},
				"period":{
					"type":"integer"
				},
				"fact":{
					"type":"integer"
				}
			}
		}
	}
}`

func elasticConnect() (context.Context, *elastic.Client, error) {
	ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetURL("http://192.168.99.100:9200", "http://192.168.99.100:9200"),
		elastic.SetSniff(false),
		elastic.SetBasicAuth("elastic", "changeme"),
	)
	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping("http://192.168.99.100:9200").Do(ctx)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://192.168.99.100:9200")
	if err != nil {
		return nil, nil, err
	}

	log.Printf("Elasticsearch version %s\n", esversion)

	// Use the IndexExists service to check if a specified index exists.
	client.DeleteIndex("tracking").Do(ctx)
	exists, err := client.IndexExists("tracking").Do(ctx)
	if err != nil {
		return nil, nil, err
	}

	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("tracking").BodyString(mapping).Do(ctx)
		if err != nil {
			return nil, nil, err
		}

		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
	return ctx, client, nil
}

func elasticIndex(ch chan int) error {
	var (
		id, requisits, reg_date, gov_choice,
		developer, trace_year, base, repeated, periodical, fact string
	)
	ctx, client, err := elasticConnect()
	if err != nil {
		return err
	}

	res, err := db.Query("select id, requisits, gov_choice, reg_date, trace_year, " +
		"developer, base, repeated, periodic, fact from track_base")

	if err != nil {
		return err
	}

	log.Print("Indexing started")
	for res.Next() {
		err := res.Scan(&id, &requisits, &reg_date, &gov_choice, &trace_year, &developer,
			&base, &repeated, &periodical, &fact)
		if err != nil {
			log.Print(err.Error(), " | ", id, "\n")
			return err
		}
		idx := Idx{id, requisits, reg_date, gov_choice,
			trace_year, developer, base, repeated, periodical, fact}
		_ = idx.writeIndex(ctx, client)
	}
	ch <- 1
	log.Print("Indexing complited!")
	return nil
}

func (idx Idx) updateIndex(id int64) error {
	var (
		index, requisits, reg_date, gov_choice,
		developer, trace_year, base, repeated, periodical, fact string
	)
	ctx, client, err := elasticConnect()

	ind, err := db.Query("select id, requisits, reg_date, gov_choice,"+
		"trace_year, developer, base, repeated, periodic, fact from track_base where id=?;", id)
	if err != nil {
		return err
	}

	for ind.Next() {
		err = ind.Scan(&index, &requisits, &reg_date, &gov_choice, &trace_year, &developer,
			&base, &repeated, &periodical, &fact)
		if err != nil {
			log.Print(err.Error())
			return err
		}
	}

	idx = Idx{index, requisits, reg_date, gov_choice,
		trace_year, developer, base, repeated, periodical, fact}
	_ = idx.writeIndex(ctx, client)

	return nil
}

func (idx Idx) writeIndex(ctx context.Context, client *elastic.Client) error {

	_, err := client.Index().
		Index("tracking").
		Type("trace").
		Id(string(idx.Id)).
		BodyJson(idx).
		Do(ctx)

	if err != nil {
		return err
	}
	return nil
}
