package main

import (

	"fmt"
)


type Event struct {
	ID int
	Name_and_requisits string
	Reg_Date string
	Government_choice string

}

func getData() []Event {
	var (
		id int
		name_and_requisits, reg_date, government_choice string
		events []Event
	)
	evt, err := db.Query("SELECT id, name_and_requisits, reg_date, government_choice from book WHERE id <= 30")
	if err != nil {
		fmt.Print(err.Error())
	}
	defer evt.Close()
	for evt.Next() {
		err = evt.Scan(&id, &name_and_requisits, &reg_date, &government_choice)
		if err != nil {
			fmt.Print(err.Error())
		}
		events = append(events, Event{id, name_and_requisits, reg_date, government_choice })
	}
	return events
}


