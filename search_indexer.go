package main

import (
	"sync"
	"log"
	"context"
	"gopkg.in/olivere/elastic.v5"
)

type workerRange struct {
	First int
	Last int
}

type workerRanges struct {
	Ids []workerRange
}

//Generating of worker ranges
func genWorkerRanges(workers int) workerRanges {
	var (
		highestId int
		ids workerRanges
	)
	row := db.QueryRow("SELECT MAX(Id) FROM track_base")
	row.Scan(&highestId)
// TODO: need improvement on highest id
	len := (highestId/100+1)*100
	rng := len/workers
	//last := len-(rng*workers)
	first := 0
	for i := 0; i < workers-1; i++{
		ids.Ids = append(ids.Ids, workerRange{first, rng*(i+1)})
		first = (rng*(i+1))+1
	}
	ids.Ids = append(ids.Ids, workerRange{(rng*(workers-1))+1, highestId})
	log.Print(ids.Ids)
	return ids
}

func indexWorker(wg *sync.WaitGroup, ctx context.Context, client *elastic.Client, f int, l int)  {

	var (
		id, requisits, reg_date, gov_choice,
		developer, trace_year, base, repeated, periodical, fact string
	)


	res, err := db.Query("select id, requisits, gov_choice, reg_date, trace_year, " +
		"developer, base, repeated, periodic, fact from track_base WHERE id BETWEEN ? AND ?", f, l)
	if err != nil {
		log.Print("ERROR:",err)
	}
	log.Print("Indexing started", f, l)
	for res.Next() {
		err := res.Scan(&id, &requisits, &reg_date, &gov_choice, &trace_year, &developer,
			&base, &repeated, &periodical, &fact)
		if err != nil {
			log.Print(err.Error(), " | ", id, "\n")
		}
		idx := Idx{id, requisits, reg_date, gov_choice,
			trace_year, developer, base, repeated, periodical, fact}

		_ = idx.writeIndex(ctx, client)
	}

	log.Print("Indexing done\n")
	wg.Done()
}


func createWorkerPool(noOfWorkers int, ch chan bool) {
	ranges := genWorkerRanges(noOfWorkers)
	var wg sync.WaitGroup

	ctx, client, err := elasticConnect()
	if err != nil {
		log.Printf("ERROR Elastic: %s\n", err)
	}
	
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go indexWorker(&wg, ctx, client, ranges.Ids[i].First, ranges.Ids[i].Last)
	}

	ch <- true
	wg.Wait()
}