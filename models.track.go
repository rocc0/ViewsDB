package main

import (
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"



)

type Trace struct {
	Fields map[string]interface{}
}

type Period struct {
	Period_id		int		`json:"pid"`
	Trace_id 		int 	`json:"trace_id"`
	Term_per 		string	`json:"term_per"`
	Res_per_bool 	int		`json:"res_per_bool"`
	Res_per_year 	int		`json:"res_per_year"`
	Res_per 		string	`json:"res_per"`
	Sign_per 		string	`json:"sign_per"`
	Publ_per 		string	`json:"publ_per"`
	Give_per 		string	`json:"give_per"`
	Concl_per 		string	`json:"concl_per"`
	Cp_bool 		int		`json:"cp_bool"`
}


func editView(name, data string, table, id int) error{
	tbl := "track_base"
	if table == 1 {
		tbl = "track_period"
	}
	log.Print("exec ",name," ", data," ",id)
	stmt, err := db.Prepare("update " + tbl + " set " + name + "= ? where id= ?;")
	if err != nil {
		log.Print(err.Error())
	} else {
		_, err := stmt.Exec(data, id)
		check(err)
	}
	return nil
}

func deleteItem(item_id int, tbl string) error {
	var table string
	if tbl == "p" {
		table = "track_period"
	} else if tbl == "b" {
		table = "track_base"
	}
	if stmt, err := db.Prepare("DELETE FROM " + table + " WHERE id=?"); err != nil {
		log.Print("\n",err,item_id,table,"\n")
		return err
	} else {
		if res, err := stmt.Exec(item_id); err != nil {
			log.Print("\n",err,res,"\n")
			return err
		}
	}
	return nil
}

func getBasicData(id int) (*Trace, error){
	var (
		trace Trace
	)
	rows, _ := db.Query("select * from track_base where id = ?;", id)
	colNames, _ := rows.Columns()

	for rows.Next() {
		columns := make([]interface{}, len(colNames))
		columnPointers := make([]interface{}, len(colNames))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}
		err := rows.Scan(columnPointers...)
		check(err)

		m := make(map[string]interface{})
		for i, colName := range colNames {
			a := columnPointers[i].(*interface{})
			if *a == nil {
				*a = []uint8("")
			}
			m[colName] =  *a
		}
		trace = Trace{m}
	}
	defer rows.Close()
	return &trace, nil
}

func getPeriodicData(id int) (*[]Period, error){
	var (
		trace_id,res_per_bool,res_per_year,cp_bool,period_id int
		term_per,sign_per,publ_per,give_per,res_per,concl_per string
		periods []Period
	)

	pers, err := db.Query("select id,trace_id,term_per,res_per_bool," +
		"COALESCE(res_per_year, '') as res_per_year," +
		"COALESCE(res_per, 'висновок відсутній') as res_per," +
		"COALESCE(sign_per, '') as res_per," +
		"COALESCE(publ_per, '') as publ_per," +
		"COALESCE(give_per, '') as give_per," +
		"COALESCE(concl_per, '') as concl_per," +
		"cp_bool from track_period where trace_id = ?;", id)
	check(err)
	for pers.Next() {
		err := pers.Scan(&trace_id,&period_id,&term_per,&res_per_bool,&res_per_year,&res_per,
			&sign_per,&publ_per,&give_per,&concl_per,&cp_bool)
		check(err)

		periods = append(periods, Period{trace_id,period_id,term_per,res_per_bool,
		res_per_year,res_per,sign_per, publ_per,give_per,concl_per,cp_bool})
	}
	defer pers.Close()
	return &periods, nil
}

func createNewItem(formData map[string]interface{}) (int, error) {

	var colNames []string
	var values []interface{}
	var phs string
	table := "track_base "
	if _, ok := formData["trace_id"]; ok {
		table = "track_period "
	}
	for k, v := range formData {
		// check that column is valid
		colNames = append(colNames, k)
		if v == "" {
			values = append(values, nil)
		} else if v == nil {
			values = append(values, " ")
		} else {
			values = append(values, v)
		}

		phs += "?,"
	}

	if len(phs) > 0 {
		phs = phs[:len(phs)-1]
	}
	phs = " values(" + phs + ");"

	colNamesString := "(" + strings.Join(colNames, ",") + ")"
	stmt, err := db.Prepare("insert into " + table + colNamesString + phs)
	check(err)

	res, err := stmt.Exec(values...)
	check(err)

	id, err := res.LastInsertId()
	check(err)

	if _, ok := formData["requisits"]; ok {
		updateIndex(id)
		return int(id), nil
	} else {
		return int(id), nil
	}

}



