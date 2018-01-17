package main

import (
	"context"
	"log"

	elastic "gopkg.in/olivere/elastic.v5"

)

// Tweet is a structure used for serializing/deserializing data in Elasticsearch.
type Idx struct {
	Id 					string		`json:"id"`
	Requisits 			string		`json:"requisits"`
	Reg_Date           	string		`json:"reg_date"`
	Gov_choice  		string		`json:"gov_choice"`
	Trace_year	 	  	string		`json:"year"`
	Developer		   	string		`json:"developer"`
	Base	 			string		`json:"base"`
	Repeated	 		string		`json:"repeat"`
	Periodical	 		string		`json:"period"`
	Fact	 			string		`json:"fact"`
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



func elasticConnect() (context.Context, *elastic.Client, error){
	ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetURL("http://192.168.99.100:9200", "http://192.168.99.100:9200"),
		elastic.SetSniff(false),
		elastic.SetBasicAuth("elastic", "changeme"),
	)
	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping("http://192.168.99.100:9200").Do(ctx)
	check(err)

	log.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://192.168.99.100:9200")
	check(err)

	log.Printf("Elasticsearch version %s\n", esversion)

	// Use the IndexExists service to check if a specified index exists.
	client.DeleteIndex("tracking").Do(ctx)
	exists, err := client.IndexExists("tracking").Do(ctx)
	check(err)

	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("tracking").BodyString(mapping).Do(ctx)
		check(err)

		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
	return ctx, client, nil
}


func elasticIndex(){
	var (
		id, requisits, reg_date, gov_choice,
		developer, trace_year, base, repeated, periodical, fact string
		trk Idx
	)
	ctx, client, err := elasticConnect()
	check(err)

	res, err := db.Query("select id, requisits, gov_choice, reg_date, trace_year, " +
		"developer, base, repeated, periodic, fact from track_base")
	check(err)

	log.Print("Indexing started")
	for res.Next(){
		err := res.Scan(&id, &requisits, &reg_date, &gov_choice, &trace_year, &developer,
			&base, &repeated, &periodical, &fact)
		if err != nil {
			log.Print(err.Error(), " | " ,id, "\n")
		}
		trk = Idx{id, requisits,reg_date,gov_choice,
			trace_year, developer, base, repeated, periodical, fact}
		_, err = client.Index().
			Index("tracking").
			Type("trace").
			Id(id).
			BodyJson(trk).
			Do(ctx)
		check(err)
	}
	log.Print("Indexing complited!")
}

func updateIndex(id int64) {
	var (
		id_ind, requisits, reg_date, gov_choice,
		developer, trace_year, base, repeated, periodical, fact string
	)
	ctx, client, err := elasticConnect()

	ind, err := db.Query("select id, requisits, reg_date, gov_choice," +
		"trace_year, developer, base, repeated, periodic, fact from track_base where id=?;", id)
	check(err)

	for ind.Next() {
		err = ind.Scan(&id_ind, &requisits, &reg_date, &gov_choice, &trace_year, &developer,
			&base, &repeated, &periodical, &fact)
		if err != nil {
			log.Print(err.Error())
		}
	}
	idx := Idx{id_ind, requisits, reg_date, gov_choice,
		trace_year,developer, base, repeated, periodical, fact}
	_, err = client.Index().
		Index("tracking").
		Type("trace").
		Id(string(id)).
		BodyJson(idx).
		Do(ctx)
}