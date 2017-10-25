package main

import (
	"log"
)

type Gov struct {
	Id int
	Name string
}


func editView(name, data string, id int) error{
	log.Print("exec ",name," ", data," ",id)
	stmt, err := db.Prepare("update views set " + name + "= ? where id= ?;")
	if err != nil {
		log.Print(err.Error())
	} else {
		_, err := stmt.Exec(data, id)
		if err != nil {
			log.Print(err.Error())
		}
	}
	return nil
}

func getReportData() (*Rating, error){
	var (
		rating [][]interface{}
		ratings Rating
		fullname string
		fullnames []string
	)
	rows, err := db.Query("SELECT gov_name, rep_total,rep_basic,rep_repited,rep_periodic,broken_period,broken_sign," +
		"broken_promo,broken_term_rep,broken_track_meth,broken_data_assump,broken_indexes,broken_my_rating,broken_dev_rating" +
		" FROM ratings INNER JOIN governments ON ratings.gov_id=governments.id;")
	if err != nil {
		log.Print(err)
	}
	colnames, _ := rows.Columns()
	names, _ := db.Query("SELECT fullname FROM rating_col_names")
	for names.Next() {
		err = names.Scan(&fullname)
		if err != nil {
			log.Print(err.Error())
		}
		fullnames = append(fullnames, fullname)
	}

	for rows.Next() {
		arr := make([]string, len(colnames))
		pointers := make([]interface{}, len(colnames))
		for i, _ := range arr {
			pointers[i] = &arr[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}
		rating = append(rating, pointers)
	}
	ratings = Rating{fullnames, rating}
	return &ratings, nil
}

func getGovs() (*[]Gov, error){
	var (
		govs []Gov
		gov_id int
		gov_name string
	)
	res, err := db.Query("SELECT id, gov_name FROM governments")
	if err != nil {
		log.Print(err.Error())
	}
	for res.Next() {
		err = res.Scan(&gov_id, &gov_name)
		if err != nil {
			log.Print(err.Error())
		}
		govs = append(govs, Gov{gov_id, gov_name })
	}
	return &govs, nil
}