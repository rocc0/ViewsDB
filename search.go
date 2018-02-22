package main

import (
	"log"
	"context"

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

//Connect to elastic
func elasticConnect() (context.Context, *elastic.Client, error) {
	ctx := context.Background()

	client, err := elastic.NewClient(
		elastic.SetURL(config.ElasticUrl, config.ElasticUrl),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(config.ElasticLog, config.ElasticPass),
	)
	if err != nil {
		return nil, nil, err
	}
	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping(config.ElasticUrl).Do(ctx)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	client.DeleteIndex("tracking").Do(ctx)
	// Use the IndexExists service to check if a specified index exists.
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

//Get item from DB and adding it to elastic
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

//Writing item to elastic index
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
