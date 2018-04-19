package main

import (
	"strings"

	"log"

	"fmt"

	"./gen"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

//BasicTrace is a part of the trace page that include register info and basic+repeated traces
type BasicTrace struct {
	Fields map[string]interface{}
}

//NewTrace is a struct that used to create new trace
type NewTrace struct {
	Info     map[string]interface{} `json:"info"`
	Basic    map[string]interface{} `json:"basic"`
	Repeated map[string]interface{} `json:"repeated"`
}

func (new NewTrace) createNewTrace() (string, error) {
	var (
		idx indexItem
	)

	traceID := gen.Generate(20)
	new.Info["trace_id"], new.Basic["trace_id"], new.Repeated["trace_id"] = traceID, traceID, traceID
	_, err := createNewSubTrace(new.Info, "trace_info")
	if err != nil {
		return "", err
	}
	_, err = createNewSubTrace(new.Basic, "trace_basic")
	if err != nil {
		return "", err
	}
	_, err = createNewSubTrace(new.Repeated, "trace_repeat")
	if err != nil {
		return "", err
	}

	if err = idx.updateIndex(traceID); err != nil {
		log.Print("2 ", err)
		return "", err
	}

	return traceID, nil
}

func (t BasicTrace) createNewPeriod() (int, error) {
	id, err := createNewSubTrace(t.Fields, "trace_period")

	if err != nil {
		return 0, err
	}

	return id, nil
}

func createNewSubTrace(data map[string]interface{}, table string) (int, error) {
	var (
		colNames []string
		values   []interface{}
		phs      string
		lastID   int
	)
	count := 1
	for k, v := range data {
		// check that column is valid
		colNames = append(colNames, k)
		if v == "" {
			values = append(values, nil)
		} else if v == nil {
			values = append(values, " ")
		} else {
			values = append(values, v)
		}
		phs += fmt.Sprintf("$%v,", count)
		count += 1
	}

	if len(phs) > 0 {
		phs = phs[:len(phs)-1]
	}
	phs = " VALUES (" + phs + ")"

	colNamesString := "(" + strings.Join(colNames, ",") + ")"
	stmt, err := db.Prepare("INSERT INTO " + table + colNamesString + phs + " RETURNING id;")
	if err != nil {
		return 0, err
	}

	err = stmt.QueryRow(values...).Scan(&lastID)
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func (t *BasicTrace) getBasicData(id string) (*BasicTrace, error) {
	rows, err := db.Query("SELECT * FROM trace_info i LEFT JOIN trace_basic b ON"+
		" i.trace_id = b.trace_id LEFT JOIN trace_repeat r ON"+
		" i.trace_id = r.trace_id WHERE i.trace_id = $1;", id)
	if err != nil {
		return nil, err
	}
	colNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	columns := make([]interface{}, len(colNames))
	columnPointers := make([]interface{}, len(colNames))
	m := make(map[string]interface{})
	for i := range columns {
		columnPointers[i] = &columns[i]
	}
	for rows.Next() {
		err := rows.Scan(columnPointers...)
		if err != nil {
			return nil, err
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
	return t, nil
}

func (t *BasicTrace) getPeriodicData(id string) (*[]map[string]interface{}, error) {
	var periods []map[string]interface{}
	pers, err := db.Query("SELECT id,trace_id, termin_zakon,termin_fact,result_bool,result_year, "+
		"result,result_comment, signed, publicated,gived,cnclsn,cnclsn_bool, cnclsn_comment,"+
		"br_my_rating,br_my_rating_c,br_dev_rating, br_dev_rating_c, p_comment "+
		"FROM trace_period WHERE trace_id = $1;", id)
	if err != nil {
		log.Print(err, "   ###1")
		return nil, err
	}
	colNames, err := pers.Columns()
	if err != nil {
		log.Print(err, "   ###2")
		return nil, err
	}
	columns := make([]interface{}, len(colNames))
	columnPointers := make([]interface{}, len(colNames))
	m := make(map[string]interface{})
	for i := range columns {
		columnPointers[i] = &columns[i]
	}
	for pers.Next() {
		err := pers.Scan(columnPointers...)
		if err != nil {
			log.Print(err, "   ###3")
			return nil, err
		}
		for i, colName := range colNames {
			a := columnPointers[i].(*interface{})
			if *a == nil {
				*a = []uint8("")
			}
			m[colName] = *a
		}
		periods = append(periods, m)
	}
	defer pers.Close()
	return &periods, nil
}

//saveTraceField saves changes maded to trace fields
func (s saveRequest) saveTraceChanges() error {
	table := map[string]string{"i": "trace_info", "b": "trace_basic", "r": "trace_repeat", "p": "trace_period"}
	field := "trace_id"
	if s.TraceType == "p" {
		field = "id"
	}

	stmt, err := db.Prepare("UPDATE " + table[s.TraceType] + " SET " + s.Name + "= $1 WHERE " + field + " = $2;")
	if err != nil {
		log.Print("1 ", err)
		return err
	}
	_, err = stmt.Exec(s.Data, s.ID)
	if err != nil {
		log.Print("2 ", err)
		return err
	}

	return nil
}

//deleteItem deletes periodic item of trace
func (d deleteRequest) deleteItem() error {
	var table string
	if d.Table == "p" {
		table = "trace_period"
	} else if d.Table == "b" {
		table = "track_basic"
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
