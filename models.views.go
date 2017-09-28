package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

type View_mainpage struct {
	Id int
	Name_and_requisits string
	Reg_Date           string
	Government_choice string
	Termin_basic string
}

type View struct {
	Fields map[string]interface{}
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(192.168.99.100:3306)/db")
	if err != nil {
		fmt.Print(err.Error())
	}

	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}
}

func getColsNames() []string {
	rows, err := db.Query("select * from book where id = 1")
	if err != nil {
		return nil
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	return cols[1:]
}

func getAllArticles() []View_mainpage {
	var (
		id int
		name_and_requisits, reg_date, government_choice, termin_basic string
		views []View_mainpage
	)
	res, err := db.Query("select id, name_and_requisits, reg_date, government_choice, termin_basic from book WHERE id <= 10")
	if err != nil {
		fmt.Print(err.Error())
	}
	defer res.Close()
	for res.Next() {
		err = res.Scan(&id, &name_and_requisits, &reg_date, &government_choice, &termin_basic)
		if err != nil {
			fmt.Print(err.Error())
		}
		views = append(views, View_mainpage{id, name_and_requisits, reg_date, government_choice, termin_basic})
	}
	return views

}

func getViewById(id int) (*View, error){
	var vie View
	rows, _ := db.Query("select * from book where id = ?", id)
	colNames, _ := rows.Columns()
	defer rows.Close()
	for rows.Next() {
		columns := make([]interface{}, len(colNames))
		columnPointers := make([]interface{}, len(colNames))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}
		m := make(map[string]interface{})
		for i, colName := range colNames {
			a := columnPointers[i].(*interface{})
			if *a == nil {
				*a = []uint8("не заповнено")
			}
			val := a
			m[colName] =  *val
		}
		vie = View{m}
	}
	return &vie, nil
}

func createNewView(formData []interface{}) (*View_mainpage, error) {

	res, err := db.Exec("insert into book (name_and_requisits, reg_date, id_reg_number, actuality_date, act_developer, government_choice, year_of_tracing, termin_basic, termin_repeated, termin_periodical," +
		"result_basic, result_repeated, result_periodical, signation_basic, signation_repeated, signation_periodical," +
			"publication_basic, publication_repeated, publication_periodical, gived_basic, gived_repeated, gived_periodical, conclusion_basic, conclusion_repeated, conclusion_periodical) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		formData...)
	if err != nil {
		log.Printf("%#v", formData)
		log.Fatal(err)
	}
	id, err := res.LastInsertId()
	vie := View_mainpage{int(id), formData[0].(string),formData[1].(string), formData[2].(string), formData[3].(string)}
	return &vie, nil
}