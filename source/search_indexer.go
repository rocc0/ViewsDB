package main

import (
	"context"
	"log"
	"os"
	"runtime/pprof"
	"sync"

	"gopkg.in/olivere/elastic.v5"
)

type (
	workerRange struct {
		First int
		Last  int
	}

	workerRanges struct {
		IDs []workerRange
	}
)

//Generating of worker ranges
func genWorkerRanges(workers int) workerRanges {
	var (
		highestID int
		ids       workerRanges
	)
	row := db.QueryRow("SELECT MAX(id) FROM trace_info")
	row.Scan(&highestID)

	// TODO: need improvement on highestID
	length := (highestID/100 + 1) * 100
	rng := length / workers
	//last := length -(rng*workers)
	first := 0
	for i := 0; i < workers-1; i++ {
		ids.IDs = append(ids.IDs, workerRange{first, rng * (i + 1)})
		first = (rng * (i + 1)) + 1
	}
	ids.IDs = append(ids.IDs, workerRange{(rng * (workers - 1)) + 1, highestID})
	return ids
}

func indexWorker(ctx context.Context, wg *sync.WaitGroup, client *elastic.Client, f int, l int) {
	var (
		id int
		traceID, regName, regDate, govChoice,
		developer, traceYear, basic, repeated, periodical, fact string
	)

	res, err := db.Query("select id, trace_id, reg_name, reg_date, gov_choice, trace_year, "+
		"developer, trace_basic, trace_repeated, trace_periodic, "+
		"trace_fact from trace_info WHERE id BETWEEN ? AND ?", f, l)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Indexing started", f, l)
	for res.Next() {
		err := res.Scan(&id, &traceID, &regName, &regDate, &govChoice, &traceYear, &developer,
			&basic, &repeated, &periodical, &fact)
		if err != nil {
			log.Fatal(err)
		}
		idx := indexItem{traceID, regName, regDate, govChoice,
			developer, traceYear, basic, repeated, periodical, fact}

		err = idx.writeIndex(ctx, client)
		if err != nil {
			log.Fatal(1, err)
		}
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
		go indexWorker(ctx, &wg, client, ranges.IDs[i].First, ranges.IDs[i].Last)
	}

	ch <- true
	wg.Wait()
}

func doReindex(workers *int, ch chan bool) {
	createWorkerPool(*workers, ch)

	if config.CPUProf != "" {
		f, err := os.Create(config.CPUProf)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		f.Close()
	}

	pprof.StopCPUProfile()
	if config.MemProf != "" {
		f, err := os.Create(config.MemProf)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
	}
	<-ch
}
