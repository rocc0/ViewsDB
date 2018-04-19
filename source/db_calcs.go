package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type governments struct {
	Gov map[string]string `json:"gov"`
}

func calculateRates(ch chan int) error {
	gvs, err := getGovernsList()
	if err != nil {
		return err
	}

	for _, v := range *gvs {
		//total trackings
		total, err := db.Query("UPDATE ratings SET rep_total=(SELECT COUNT(developer) "+
			"FROM track_base WHERE developer=? AND trace_year=?) WHERE gov_id=?", v.ID, "2016", v.ID)
		if err != nil {
			return err
		}
		total.Close()
		//basic trackings
		base, _ := db.Query("UPDATE ratings SET rep_basic=(SELECT COUNT(term_basic) "+
			"FROM track_base WHERE trace_year=? AND developer=? AND track_type=?) WHERE gov_id=?",
			"2016", v.ID, "base", v.ID)
		base.Close()

		//repeated trackings
		rep, _ := db.Query("UPDATE ratings SET rep_repited=(SELECT COUNT(termin_rep) "+
			"FROM track_base WHERE trace_year=? AND developer=? AND track_type=?) WHERE gov_id=?",
			"2016", v.ID, "repeated", v.ID)
		rep.Close()
		//periodic trackings
		periodic, err := db.Query("UPDATE ratings SET rep_periodic=(SELECT COUNT(developer) "+
			"FROM track_base WHERE trace_year=? AND developer=? AND track_type=?) WHERE gov_id=?",
			"2016", v.ID, "periodic", v.ID)
		if err != nil {
			return err
		}
		periodic.Close()

		//broken
		brokenPeriod, _ := db.Query("UPDATE ratings SET broken_period=(SELECT COUNT(broken_period) "+
			"FROM track_base WHERE trace_year=? AND broken_period = 1 AND developer=?) WHERE gov_id=?",
			"2016", v.ID, v.ID)
		brokenPeriod.Close()

		brokenSign, _ := db.Query("UPDATE ratings SET broken_sign=(SELECT COUNT(broken_sign) "+
			"FROM track_base WHERE trace_year=?  AND broken_sign = 1 AND developer=?) WHERE gov_id=?",
			"2016", v.ID, v.ID)
		brokenSign.Close()

		brokenPromo, _ := db.Query("UPDATE ratings SET broken_promo=(SELECT COUNT(broken_promo) "+
			"FROM track_base WHERE trace_year=? AND broken_promo = 1 AND developer=?) WHERE gov_id=?",
			"2016", v.ID, v.ID)
		brokenPromo.Close()

		brokenTermRep, _ := db.Query("UPDATE ratings SET broken_term_rep=(SELECT COUNT(broken_term_rep) "+
			"FROM track_base WHERE trace_year=? AND broken_term_rep = 1 AND developer=?) WHERE gov_id=?",
			"2016", v.ID, v.ID)
		brokenTermRep.Close()

		brokenTrackMeth, _ := db.Query("UPDATE ratings SET broken_track_meth=(SELECT COUNT(broken_track_meth) "+
			"FROM track_base WHERE trace_year=?  AND broken_track_meth = 1 AND developer=?) WHERE gov_id=?",
			"2016", v.ID, v.ID)
		brokenTrackMeth.Close()

		brokenDataAssump, _ := db.Query("UPDATE ratings SET broken_data_assump=(SELECT COUNT(broken_data_assump) "+
			"FROM track_base WHERE trace_year=?  AND broken_data_assump = 1 AND developer=?) WHERE gov_id=?",
			"2016", v.ID, v.ID)
		brokenDataAssump.Close()

		brokenIndexes, _ := db.Query("UPDATE ratings SET broken_indexes=(SELECT COUNT(broken_indexes) "+
			"FROM track_base WHERE trace_year=? AND broken_indexes = 1 AND developer=?) WHERE gov_id=?",
			"2016", v.ID, v.ID)
		brokenIndexes.Close()

		//broken period

		brokenMyRating, _ := db.Query("UPDATE ratings SET broken_my_rating=(SELECT COUNT(broken_my_rating) "+
			"FROM track_period WHERE term_per between ? AND ? AND broken_my_rating = 1 AND developer=?) WHERE gov_id=?",
			"2016-01-01", "2016-12-31", v.ID, v.ID)

		brokenMyRating.Close()

		brokenDevRating, _ := db.Query("UPDATE ratings SET broken_dev_rating=(SELECT COUNT(broken_dev_rating) "+
			"FROM track_period WHERE term_per between ? AND ? AND broken_dev_rating = 1 AND developer=?) WHERE gov_id=?",
			"2016-01-01", "2016-12-31", v.ID, v.ID)
		brokenDevRating.Close()
	}
	ch <- 1

	return nil
}

func calulateRatesSecond() error {
	gvs, err := getGovernsList()
	if err != nil {
		return err
	}
	for _, v := range *gvs {
		//total trackings
		total, _ := db.Query("UPDATE ratings2 SET rep_total=(SELECT COUNT(developer) "+
			"FROM track_base WHERE developer=? AND ) WHERE gov_id=?", v.ID, v.ID)
		total.Close()
	}
	return nil
}

func getReportData() (*[]string, *[]governments, error) {
	var (
		rating  governments
		ratings []governments
		column  string
		columns []string
	)
	rows, err := db.Query("SELECT gov_name, rep_total,rep_basic,rep_repited,rep_periodic,broken_period,broken_sign," +
		"broken_promo,broken_term_rep,broken_track_meth,broken_data_assump,broken_indexes,broken_my_rating,broken_dev_rating" +
		" FROM ratings INNER JOIN governments ON ratings.gov_id=governments.id;")
	if err != nil {
		return nil, nil, err
	}
	colNames, _ := rows.Columns()
	names, _ := db.Query("SELECT fullname FROM rating_cols")
	for names.Next() {
		err = names.Scan(&column)
		if err != nil {
			return nil, nil, err
		}

		columns = append(columns, column)
	}

	for rows.Next() {
		columns = make([]string, len(colNames))
		columnPointers := make([]interface{}, len(colNames))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		err := rows.Scan(columnPointers...)
		if err != nil {
			return nil, nil, err
		}

		m := make(map[string]string)
		for i, colName := range colNames {
			a := columnPointers[i].(*string)
			m[colName] = *a
		}
		rating = governments{m}
		ratings = append(ratings, rating)
	}

	return &columns, &ratings, nil
}
