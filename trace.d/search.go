package main

import (
	"context"

	"gopkg.in/olivere/elastic.v5"
)

type indexItem struct {
	TraceID    string `json:"trace_id"`
	RegName    string `json:"reg_name"`
	RegDate    string `json:"reg_date"`
	GovChoice  string `json:"gov_choice"`
	Developer  string `json:"developer"`
	TraceYear  string `json:"trace_year"`
	Basic      string `json:"trace_basic"`
	Repeated   string `json:"trace_repeat"`
	Periodical string `json:"trace_period"`
	Fact       string `json:"trace_fact"`
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
				"trace_id":{
					"type":"string"
				},
				"reg_name":{
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
				"trace_year":{
					"type":"integer"
				},
				"trace_basic":{
					"type":"integer"
				},
				"trace_repeat":{
					"type":"integer"
				},
				"trace_period":{
					"type":"integer"
				},
				"trace_fact":{
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
		elastic.SetURL(config.ElasticURL, config.ElasticURL),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(config.ElasticLog, config.ElasticPass),
	)
	if err != nil {
		return nil, nil, err
	}

	//client.DeleteIndex("tracking").Do(ctx)

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

//Getting item from DB and adding it to elastic
func (idx indexItem) updateIndex(id string) error {
	var index int
	ctx, client, err := elasticConnect()
	if err != nil {
		return err
	}

	stmt := db.QueryRow("select id, trace_id, reg_name, reg_date, gov_choice,"+
		"trace_year, developer, trace_basic, trace_repeated, trace_periodic,"+
		" trace_fact from trace_info where trace_id=$1;", id)

	if err = stmt.Scan(&index, &idx.TraceID, &idx.RegName, &idx.RegDate, &idx.GovChoice,
		&idx.TraceYear, &idx.Developer, &idx.Basic, &idx.Repeated, &idx.Periodical, &idx.Fact); err != nil {
		return err
	}

	if err = idx.writeIndex(ctx, client); err != nil {
		return err
	}
	return nil
}

//Writing item to elastic index
func (idx indexItem) writeIndex(ctx context.Context, client *elastic.Client) error {

	if _, err := client.Index().
		Index("tracking").
		Type("trace").
		Id(idx.TraceID).
		BodyJson(idx).
		Do(ctx); err != nil {
		return err
	}
	return nil
}
