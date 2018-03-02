package main

import (
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//BasicTrace is a part of the trace page that include register info and basic+repeated traces
type BasicTrace struct {
	Fields map[string]interface{}
}

//PeriodTrace is a part of the trace page thar include only periodic traces
type PeriodTrace struct {
	PeriodID               int    `json:"pid"`
	TraceID                int    `json:"trace_id"`
	TermPer                string `json:"term_per"`
	ResPerBool             int    `json:"res_per_bool"`
	ResPerYear             int    `json:"res_per_year"`
	ResPerComment          string `json:"res_per_comment"`
	ResPer                 string `json:"res_per"`
	SignPer                string `json:"sign_per"`
	PublPer                string `json:"publ_per"`
	GivePer                string `json:"give_per"`
	ConclPer               string `json:"concl_per"`
	ConclPerBool           int    `json:"cp_bool"`
	ConclPerComment        string `json:"concl_per_comment"`
	BrokenMyRating         int    `json:"broken_my_rating"`
	BrokenMyRatingComment  string `json:"broken_my_rating_c"`
	BrokenDevRating        int    `json:"broken_dev_rating"`
	BrokenDevRatingComment string `json:"broken_dev_rating_c"`
}

func (s saveRequest) saveTraceField() error {
	table := "track_base"
	if s.TraceType == 1 {
		table = "track_period"
	}
	stmt, err := db.Prepare("UPDATE " + table + " SET " + s.Name + "= ? WHERE id= ?;")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(s.Data, s.ID)
	if err != nil {
		return err
	}

	return nil
}

func (d deleteRequest) deleteItem() error {
	var table string
	if d.TableName == "p" {
		table = "track_period"
	} else if d.TableName == "b" {
		table = "track_base"
	}
	stmt, err := db.Prepare("DELETE FROM " + table + " WHERE id=?")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(d.TraceID); err != nil {
		return err
	}
	return nil
}

func (t *BasicTrace) getBasicData(id int) error {
	rows, _ := db.Query("select * from track_base where id = ?;", id)
	colNames, _ := rows.Columns()
	columns := make([]interface{}, len(colNames))
	columnPointers := make([]interface{}, len(colNames))
	m := make(map[string]interface{})
	for i := range columns {
		columnPointers[i] = &columns[i]
	}
	for rows.Next() {
		err := rows.Scan(columnPointers...)
		if err != nil {
			return err
		}
		for i, colName := range colNames {
			a := columnPointers[i].(*interface{})
			if *a == nil {
				*a = []uint8("")
			}
			m[colName] = *a
		}
		t.Fields = m
	}
	defer rows.Close()
	return nil
}

func getPeriodicData(id int) (*[]PeriodTrace, error) {
	var (
		traceID, resPerBool, resPerYear, periodID, ConclPerBool, brokenMyRating, brokenDevRating int
		termPer, resPer, resPerComment, signPer, publPer, givePer,
		conclPer, conclPerComment, brokenMyRatingComment, brokenDevRatingComment string
		periods []PeriodTrace
	)

	pers, err := db.Query("select id,trace_id,term_per,res_per_bool, "+
		"COALESCE(res_per_year, 0) as res_per_year, "+
		"COALESCE(res_per, 'висновок відсутній') as res_per,"+
		"COALESCE(res_per_comment, 'коментар відсутній') as res_per_comment, "+
		"COALESCE(sign_per, '') as res_per,"+
		"COALESCE(publ_per, '') as publ_per,"+
		"COALESCE(give_per, '') as give_per,"+
		"COALESCE(concl_per, '') as concl_per,"+
		"cp_bool,"+
		"COALESCE(concl_per_comment, 'коментар відсутній') as concl_per_comment,"+
		"COALESCE(broken_my_rating, 0) as broken_my_rating,"+
		"COALESCE(broken_my_rating_c, '') as broken_my_rating_c,"+
		"COALESCE(broken_dev_rating, 0) as broken_dev_rating, "+
		"COALESCE(broken_my_rating_c, '') as broken_dev_rating_c "+
		"from track_period where trace_id = ?;", id)
	if err != nil {
		return nil, err
	}
	for pers.Next() {
		err := pers.Scan(&traceID, &periodID, &termPer, &resPerBool, &resPerYear, &resPer, &resPerComment,
			&signPer, &publPer, &givePer, &conclPer, &ConclPerBool, &conclPerComment, &brokenMyRating, &brokenMyRatingComment,
			&brokenDevRating, &brokenDevRatingComment)
		if err != nil {
			return nil, err
		}

		periods = append(periods, PeriodTrace{traceID, periodID, termPer, resPerBool,
			resPerYear, resPerComment, resPer, signPer, publPer,
			givePer, conclPer, ConclPerBool, conclPerComment,
			brokenMyRating, brokenMyRatingComment,
			brokenDevRating, brokenDevRatingComment})
	}
	defer pers.Close()
	return &periods, nil
}

func (t BasicTrace) createNewItem() (int, error) {
	var (
		colNames []string
		values   []interface{}
		phs      string
		idx      indexItem
	)
	table := "track_base "

	if _, ok := t.Fields["trace_id"]; ok {
		table = "track_period "
	}
	for k, v := range t.Fields {
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
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(values...)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if _, ok := t.Fields["requisits"]; ok {
		idx.updateIndex(id)
		return int(id), nil
	}
	return int(id), nil

}

func (e editGovernName) editGovName() error {
	stmt, err := db.Prepare("UPDATE governments SET gov_name=? WHERE id=?;")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(e.Name, e.ID)
	if err != nil {
		return err
	}
	return nil
}
