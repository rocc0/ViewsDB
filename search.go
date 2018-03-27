package main

import (
	"context"
	"log"

	elastic "gopkg.in/olivere/elastic.v5"
)

type indexItem struct {
	ID         string `json:"id"`
	Requisits  string `json:"requisits"`
	RegDate    string `json:"reg_date"`
	GovChoice  string `json:"gov_choice"`
	TraceYear  string `json:"year"`
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

//Getting item from DB and adding it to elastic
func (idx indexItem) updateIndex(id int64) error {
	var (
		index, requisits, regDate, govChoice,
		developer, traceYear, base, repeated, periodical, fact string
	)
	ctx, client, err := elasticConnect()

	ind, err := db.Query("select id, reg_name, reg_date, gov_choice,"+
		"trace_year, developer, trace_basic, trace_repeated, trace_periodic,"+
		" trace_fact from trace_info where id=?;", id)
	if err != nil {
		return err
	}

	for ind.Next() {
		err = ind.Scan(&index, &requisits, &regDate, &govChoice, &traceYear, &developer,
			&base, &repeated, &periodical, &fact)
		if err != nil {
			log.Print(err.Error())
			return err
		}
	}

	idx = indexItem{index, requisits, regDate, govChoice,
		traceYear, developer, base, repeated, periodical, fact}
	_ = idx.writeIndex(ctx, client)

	return nil
}

//Writing item to elastic index
func (idx indexItem) writeIndex(ctx context.Context, client *elastic.Client) error {

	_, err := client.Index().
		Index("tracking").
		Type("trace").
		Id(string(idx.ID)).
		BodyJson(idx).
		Do(ctx)

	if err != nil {
		return err
	}
	return nil
}
