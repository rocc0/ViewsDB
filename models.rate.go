package main


type Column struct {
	Header string `json:"header"`
}

type Rating struct {
	Gov map[string]string `json:"gov"`
}

type Gov struct {
	Id int `json:"id"`
	Name string	`json:"name"`
}


func calculateRates()  {
	gvs, err := getGovs()
	check(err)

	for _, v := range *gvs {
		//total trackings
		total, _ := db.Query("UPDATE ratings SET rep_total=(SELECT COUNT(developer) " +
			"FROM track_base WHERE developer=?) WHERE gov_id=?",v.Id,v.Id)
		total.Close()
		//basic trackings
		base, _ :=db.Query("UPDATE ratings SET rep_basic=(SELECT COUNT(term_basic) " +
			"FROM track_base WHERE term_basic between ? and ? AND developer=?) WHERE gov_id=?",
				"2016-01-01","2016-12-31",v.Id,v.Id)
		base.Close()

		//repeated trackings
		rep, _ := db.Query("UPDATE ratings SET rep_repited=(SELECT COUNT(termin_rep) " +
			"FROM track_base WHERE termin_rep between ? and ? AND developer=?) WHERE gov_id=?",
				"2016-01-01","2016-12-31",v.Id,v.Id)
		rep.Close()
		//periodic trackings
		periodic, _ := db.Query("UPDATE ratings SET rep_periodic=(SELECT COUNT(term_per) " +
			"FROM track_period WHERE term_per between ? and ? AND developer=?) WHERE gov_id=?",
				"2016-01-01","2016-12-31",v.Id,v.Id)
		periodic.Close()
		base_term, _ :=db.Query("UPDATE ratings SET rep_basic=(SELECT COUNT(term_basic) " +
			"FROM track_base WHERE term_basic between ? and ? AND developer=?) WHERE gov_id=?",
			"2016-01-01","2016-12-31",v.Id,v.Id)
		base_term.Close()


	}

}

func getGovs() (*[]Gov, error){
	var (
		govs []Gov
		gov_id int
		gov_name string
	)
	res, err := db.Query("SELECT id, gov_name FROM governments")
	check(err)

	for res.Next() {
		err = res.Scan(&gov_id, &gov_name)
		check(err)

		govs = append(govs, Gov{gov_id, gov_name })
	}
	return &govs, nil
}

func getReportData() (*[]string,*[]Rating, error){
	var (
		rating Rating
		ratings []Rating

		column string
		columns []string
	)
	rows, err := db.Query("SELECT gov_name, rep_total,rep_basic,rep_repited,rep_periodic,broken_period,broken_sign," +
		"broken_promo,broken_term_rep,broken_track_meth,broken_data_assump,broken_indexes,broken_my_rating,broken_dev_rating" +
		" FROM ratings INNER JOIN governments ON ratings.gov_id=governments.id;")
	check(err)

	colNames, _ := rows.Columns()
	names, _ := db.Query("SELECT fullname FROM rating_cols")
	for names.Next() {
		err = names.Scan(&column)
		check(err)

		columns = append(columns, column)
	}

	for rows.Next() {
		columns := make([]string, len(colNames))
		columnPointers := make([]interface{}, len(colNames))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}
		err := rows.Scan(columnPointers...)
		check(err)

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

