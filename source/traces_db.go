package main

import (
	"strings"

	_ "github.com/go-sql-driver/mysql"
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

//PeriodTrace is a part of the trace page thar include only periodic traces
type PeriodTrace struct {
	PeriodID               int    `json:"pid"`
	TraceID                string `json:"trace_id"`
	TermZakon              string `json:"termin_zakon"`
	TermFact               string `json:"termin_fact"`
	Result                 string `json:"result"`
	ResultBool             int    `json:"result_bool"`
	ResultYear             int    `json:"result_year"`
	ResultComment          string `json:"result_comment"`
	Signed                 string `json:"signed"`
	Publicated             string `json:"publicated"`
	Gived                  string `json:"gived"`
	Cnclsn                 string `json:"cnclsn"`
	CnclsnBool             int    `json:"cnclsn_bool"`
	CnclsnComment          string `json:"cnclsn_comment"`
	BrokenMyRating         int    `json:"br_my_rating"`
	BrokenMyRatingComment  string `json:"br_my_rating_c"`
	BrokenDevRating        int    `json:"br_dev_rating"`
	BrokenDevRatingComment string `json:"br_dev_rating_c"`
	Comment                string `json:"p_comment"`
}

func (new NewTrace) createNewTrace() (string, error) {
	var (
		idx indexItem
	)

	traceID := generate(20)
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

	idx.updateIndex(traceID)

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
	)

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

		phs += "?,"
	}

	if len(phs) > 0 {
		phs = phs[:len(phs)-1]
	}
	phs = " VALUES(" + phs + ");"

	colNamesString := "(" + strings.Join(colNames, ",") + ")"
	stmt, err := db.Prepare("INSERT INTO " + table + colNamesString + phs)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(values...)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastID), nil
}

func (t *BasicTrace) getBasicData(id string) error {
	rows, err := db.Query("SELECT * FROM trace_info i LEFT JOIN trace_basic b ON"+
		" i.trace_id = b.trace_id LEFT JOIN trace_repeat r ON"+
		" i.trace_id = r.trace_id WHERE i.trace_id = ?;", id)
	if err != nil {
		return err
	}
	colNames, err := rows.Columns()
	if err != nil {
		return err
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

func getPeriodicData(id string) (*[]PeriodTrace, error) {
	var (
		resultBool, resultYear, periodID, cnclsnBool, brokenMyRating, brokenDevRating int
		traceID, termZakon, termFact, result, resultComment, signed, publicated, gived,
		cnclsn, cnclsnComment, brokenMyRatingComment, brokenDevRatingComment, comment string
		periods []PeriodTrace
	)

	pers, err := db.Query("SELECT id,trace_id,"+
		"COALESCE(termin_zakon, '') AS termin_zakon,"+
		"COALESCE(termin_fact, '') AS termin_fact,"+
		"result_bool, "+
		"COALESCE(result_year, 0) AS result_year, "+
		"COALESCE(result, 'висновок відсутній') AS result,"+
		"COALESCE(result_comment, 'коментар відсутній') AS result_comment, "+
		"COALESCE(signed, '') AS signed,"+
		"COALESCE(publicated, '') AS publicated,"+
		"COALESCE(gived, '') AS gived,"+
		"COALESCE(cnclsn, '') AS cnclsn,"+
		"cnclsn_bool, "+
		"COALESCE(cnclsn_comment, 'коментар відсутній') AS cnclsn_comment,"+
		"COALESCE(br_my_rating, 0) AS br_my_rating,"+
		"COALESCE(br_my_rating_c, '') AS br_my_rating_c,"+
		"COALESCE(br_dev_rating, 0) AS br_dev_rating, "+
		"COALESCE(br_dev_rating_c, '') AS br_dev_rating_c, "+
		"COALESCE(p_comment, '') AS p_comment "+
		"FROM trace_period WHERE trace_id = ?;", id)
	if err != nil {
		return nil, err
	}
	for pers.Next() {
		err := pers.Scan(&periodID, &traceID, &termZakon, &termFact,
			&resultBool, &resultYear, &result, &resultComment,
			&signed, &publicated, &gived, &cnclsn, &cnclsnBool,
			&cnclsnComment, &brokenMyRating, &brokenMyRatingComment,
			&brokenDevRating, &brokenDevRatingComment, &comment)
		if err != nil {
			return nil, err
		}

		periods = append(periods, PeriodTrace{periodID, traceID, termZakon, termFact,
			result, resultBool, resultYear, resultComment, signed,
			publicated, gived, cnclsn, cnclsnBool,
			cnclsnComment, brokenMyRating,
			brokenMyRatingComment, brokenDevRating,
			brokenDevRatingComment, comment})
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

	stmt, err := db.Prepare("UPDATE " + table[s.TraceType] + " SET " + s.Name + "= ? WHERE " + field + " = ?;")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(s.Data, s.ID)
	if err != nil {
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
