package main

import (
	"context"
	"fmt"
	elastic "gopkg.in/olivere/elastic.v5"
)

// Tweet is a structure used for serializing/deserializing data in Elasticsearch.
type Idx struct {
	Id string						`json:"id"`
	Name_and_requisits string		`json:"name_and_requisits"`
	Reg_Date           string		`json:"reg_date"`
	Government_choice  string		`json:"government_choice"`
	Year_of_tracing	   string		`json:"year_of_tracing"`
	Act_developer	   string		`json:"act_developer"`
}

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"view":{
			"properties":{
				"id":{
					"type":"integer"
				},
				"name_and_requisits":{
					"type":"string",
					"analyzer": "ukrainian"
				},
				"reg_date":{
					"type":"text"
				},
				"government_choice":{
					"type":"string",
					"analyzer": "ukrainian"
				},
				"year":{
					"type":"integer",
				},
				"act_developer":{
					"type":"integer",
				}
			}
		}
	}
}`

func elasticIndex(){
	ctx := context.Background()

	client, err := elastic.NewClient(
		elastic.SetURL("http://192.168.99.100:9200", "http://192.168.99.100:9200"),
		elastic.SetSniff(false),
		elastic.SetBasicAuth("elastic", "changeme"),
	)

	if err != nil {
		// Handle error
		panic(err)
	}

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping("http://192.168.99.100:9200").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://192.168.99.100:9200")
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("views").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("views").BodyString(mapping).Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
	var (
		id, name_and_requisits, reg_date, government_choice, act_developer, year_of_tracing string
		vie Idx
	)
	res, err := db.Query("select id, name_and_requisits, reg_date, government_choice, year_of_tracing, act_developer from views")
	for res.Next(){
		err := res.Scan(&id, &name_and_requisits, &reg_date, &government_choice, &year_of_tracing, &act_developer )
		if err != nil {
			fmt.Print(err.Error())
		}
		vie = Idx{id, name_and_requisits,reg_date,government_choice,
		year_of_tracing, act_developer}
		_, err = client.Index().
			Index("views").
			Type("view").
			Id(id).
			BodyJson(vie).
			Do(ctx)
		if err != nil {
			fmt.Print(err.Error())
		}
	}

}