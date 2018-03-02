package main

import (
	"context"
	"log"
	"os"
	"runtime/pprof"
	"sync"

	"gopkg.in/olivere/elastic.v5"
)

type workerRange struct {
	First int
	Last  int
}

type workerRanges struct {
	IDs []workerRange
}

//Generating of worker ranges
func genWorkerRanges(workers int) workerRanges {
	var (
		highestID int
		ids       workerRanges
	)
	row := db.QueryRow("SELECT MAX(Id) FROM track_base")
	row.Scan(&highestID)

	// TODO: need improvement on highest id
	len := (highestID/100 + 1) * 100
	rng := len / workers
	//last := len-(rng*workers)
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
		id, requisits, regDate, govChoice,
		developer, traceYear, base, repeated, periodical, fact string
	)

	res, err := db.Query("select id, requisits, gov_choice, reg_date, trace_year, "+
		"developer, base, repeated, periodic, fact from track_base WHERE id BETWEEN ? AND ?", f, l)
	if err != nil {
		log.Print("ERROR:", err)
	}
	log.Print("Indexing started", f, l)
	for res.Next() {
		err := res.Scan(&id, &requisits, &regDate, &govChoice, &traceYear, &developer,
			&base, &repeated, &periodical, &fact)
		if err != nil {
			log.Print(err.Error(), " | ", id, "\n")
		}
		idx := indexItem{id, requisits, regDate, govChoice,
			traceYear, developer, base, repeated, periodical, fact}

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
