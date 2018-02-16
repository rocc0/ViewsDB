package main

type Rating struct {
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
			"FROM track_base WHERE developer=? AND trace_year=?) WHERE gov_id=?", v.Id, "2016", v.Id)
		if err != nil {
			return err
		}
		total.Close()
		//basic trackings
		base, _ := db.Query("UPDATE ratings SET rep_basic=(SELECT COUNT(term_basic) "+
			"FROM track_base WHERE trace_year=? AND developer=? AND track_type=?) WHERE gov_id=?",
			"2016", v.Id, "base", v.Id)
		base.Close()

		//repeated trackings
		rep, _ := db.Query("UPDATE ratings SET rep_repited=(SELECT COUNT(termin_rep) "+
			"FROM track_base WHERE trace_year=? AND developer=? AND track_type=?) WHERE gov_id=?",
			"2016", v.Id, "repeated", v.Id)
		rep.Close()
		//periodic trackings
		periodic, err := db.Query("UPDATE ratings SET rep_periodic=(SELECT COUNT(developer) "+
			"FROM track_base WHERE trace_year=? AND developer=? AND track_type=?) WHERE gov_id=?",
			"2016", v.Id, "periodic", v.Id)
		if err != nil {
			return err
		}
		periodic.Close()

		//broken
		broken_period, _ := db.Query("UPDATE ratings SET broken_period=(SELECT COUNT(broken_period) "+
			"FROM track_base WHERE trace_year=? AND broken_period = 1 AND developer=?) WHERE gov_id=?",
			"2016", v.Id, v.Id)
		broken_period.Close()

		broken_sign, _ := db.Query("UPDATE ratings SET broken_sign=(SELECT COUNT(broken_sign) "+
			"FROM track_base WHERE trace_year=?  AND broken_sign = 1 AND developer=?) WHERE gov_id=?",
			"2016", v.Id, v.Id)
		broken_sign.Close()

		broken_promo, _ := db.Query("UPDATE ratings SET broken_promo=(SELECT COUNT(broken_promo) "+
			"FROM track_base WHERE trace_year=? AND broken_promo = 1 AND developer=?) WHERE gov_id=?",
			"2016", v.Id, v.Id)
		broken_promo.Close()

		broken_term_rep, _ := db.Query("UPDATE ratings SET broken_term_rep=(SELECT COUNT(broken_term_rep) "+
			"FROM track_base WHERE trace_year=? AND broken_term_rep = 1 AND developer=?) WHERE gov_id=?",
			"2016", v.Id, v.Id)
		broken_term_rep.Close()

		broken_track_meth, _ := db.Query("UPDATE ratings SET broken_track_meth=(SELECT COUNT(broken_track_meth) "+
			"FROM track_base WHERE trace_year=?  AND broken_track_meth = 1 AND developer=?) WHERE gov_id=?",
			"2016", v.Id, v.Id)
		broken_track_meth.Close()

		broken_data_assump, _ := db.Query("UPDATE ratings SET broken_data_assump=(SELECT COUNT(broken_data_assump) "+
			"FROM track_base WHERE trace_year=?  AND broken_data_assump = 1 AND developer=?) WHERE gov_id=?",
			"2016", v.Id, v.Id)
		broken_data_assump.Close()

		broken_indexes, _ := db.Query("UPDATE ratings SET broken_indexes=(SELECT COUNT(broken_indexes) "+
			"FROM track_base WHERE trace_year=? AND broken_indexes = 1 AND developer=?) WHERE gov_id=?",
			"2016", v.Id, v.Id)
		broken_indexes.Close()

		//broken period

		broken_my_rating, _ := db.Query("UPDATE ratings SET broken_my_rating=(SELECT COUNT(broken_my_rating) "+
			"FROM track_period WHERE term_per between ? AND ? AND broken_my_rating = 1 AND developer=?) WHERE gov_id=?",
			"2016-01-01", "2016-12-31", v.Id, v.Id)

		broken_my_rating.Close()

		broken_dev_rating, _ := db.Query("UPDATE ratings SET broken_dev_rating=(SELECT COUNT(broken_dev_rating) "+
			"FROM track_period WHERE term_per between ? AND ? AND broken_dev_rating = 1 AND developer=?) WHERE gov_id=?",
			"2016-01-01", "2016-12-31", v.Id, v.Id)
		broken_dev_rating.Close()
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
			"FROM track_base WHERE developer=? AND ) WHERE gov_id=?", v.Id, v.Id)
		total.Close()
	}
	return nil
}



func getReportData() (*[]string, *[]Rating, error) {
	var (
		rating  Rating
		ratings []Rating

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
		columns := make([]string, len(colNames))
		columnPointers := make([]interface{}, len(colNames))
		for i, _ := range columns {
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
		rating = Rating{m}
		ratings = append(ratings, rating)
	}

	return &columns, &ratings, nil
}
