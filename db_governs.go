package main

type Government struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
func getGovernsList() (*[]Government, error) {
	var (
		govs     []Government
		gov_id   int
		gov_name string
	)
	res, err := db.Query("SELECT id, gov_name FROM governments")
	if err != nil {
		return nil, err
	}

	for res.Next() {
		err = res.Scan(&gov_id, &gov_name)
		if err != nil {
			return nil, err
		}

		govs = append(govs, Government{gov_id, gov_name})
	}
	return &govs, nil
}
