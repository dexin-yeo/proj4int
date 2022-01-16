package update

import (
	"log"
	"notion/api"
	"strconv"
)

/*
Updates row/s in Notion
Param:
	- args: args from command line
Returns:
	- error
*/
func Update(args []string) error {
	// checks whether csv is given
	if len(args) < 1 {
		log.Fatal("not enough arguments given")
	}

	var err error
	obj := api.CreateAPI()

	// csv file given to update
	if len(args) <= 1 {
		// reads and extract data from the csv the csv file
		data := api.ReadCSVFile(args[0])

		log.Printf("Updating rows to database as pages...\n")

		// first element contains the col names
		cols := data[0]
		// loop thru the rows
		for _, v := range data[1:] {
			// remake rows
			rows := make([]interface{}, (len(v)-1)*3)
			// count for rows
			j := 0
			id := ""
			for i, c := range cols {
				if c == "page_id" {
					id = v[i]
				} else {
					var val interface{}
					val = v[i]
					if c == "table_name" || c == "is_nullable" {
						rows[j] = api.PROPERTIES_SELECT
					} else if c == "column_name" {
						rows[j] = api.PROPERTIES_TITLE
					} else if c == "ordinal_position" {
						rows[j] = api.PROPERTIES_NUMBER
						val, err = strconv.Atoi(v[i])
						if err != nil {
							return err
						}
					} else if c == "character_maximum_length" {
						rows[j] = api.PROPERTIES_NUMBER
						if val == "" {
							val = 0
						} else {
							val, err = strconv.Atoi(v[i])
							if err != nil {
								return err
							}
						}
					} else if c == "column_default" || c == "data_type" {
						rows[j] = api.PROPERTIES_TEXT
					}
					j++
					rows[j] = c
					j++
					rows[j] = val
					j++
				}
			}

			// hardcoded sequence based on csv
			// each row is a page
			// format type, colname, input
			pg, err := obj.UpdatePage(id, rows...)
			if err != nil {
				return err
			}
			log.Printf("Page %s updated\n", pg)
		}
	}

	return nil
}
