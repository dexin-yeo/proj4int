package insert

import (
	"log"
	"notion/api"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
Insert new row/s into notion
Param:
	- args: args from command line
Returns:
	- error
*/
func Insert(args []string) error {
	// checks whether database id is given
	if len(args) < 1 {
		log.Fatal("database id is not given")
	}

	// get current wd
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	s := strings.Split(wd, "/")
	// change to data path
	s[len(s)-1] = "data"
	// form back the path storing the json examples
	jsonPath := strings.Join(s, "/")

	obj := api.CreateAPI()

	// csv file given, populate data into notion
	if len(args) < 3 {
		// reads and extract data from the csv the csv file
		data := api.ReadCSVFile(args[1])

		// extra symbols
		re, err := regexp.Compile(`[^\w]`)
		if err != nil {
			return err
		}

		log.Printf("Populating rows to database %s as pages...\n", args[0])

		// first element contains the col names
		cols := data[0]
		// remove extra symbols from title name to prevent errs
		for i, v := range cols {
			cols[i] = re.ReplaceAllString(v, "")
		}

		// loop thru the rows
		for _, v := range data[1:] {
			// convert ordinal position to int
			v2, err := strconv.Atoi(v[2])
			if err != nil {
				return err
			}

			// maximum length set empty to 0
			v6 := 0
			if v[6] != "" {
				v6, err = strconv.Atoi(v[6])
				if err != nil {
					return err
				}
			}

			// hardcoded sequence based on csv
			// each row is a page
			// format type, colname, input
			pg, err := obj.CreatePage(jsonPath+"/createPage.json", args[0],
				api.PROPERTIES_SELECT, cols[0], v[0],
				api.PROPERTIES_TITLE, cols[1], v[1],
				api.PROPERTIES_NUMBER, cols[2], v2,
				api.PROPERTIES_TEXT, cols[3], v[3],
				api.PROPERTIES_SELECT, cols[4], v[4],
				api.PROPERTIES_TEXT, cols[5], v[5],
				api.PROPERTIES_NUMBER, cols[6], v6)
			if err != nil {
				return err
			}
			log.Printf("Page %s populated\n", pg)
		}
	} else if len(args) < 9 {
		// row is input along with command line
		log.Printf("Inserting row to database %s as page...\n", args[0])
		// convert ordinal position to int
		v3, err := strconv.Atoi(args[3])
		if err != nil {
			return err
		}

		// maximum length set empty to 0
		v7 := 0
		if len(args) == 8 && args[7] != "" {
			v7, err = strconv.Atoi(args[7])
			if err != nil {
				return err
			}
		}

		// hardcoded sequence
		pg, err := obj.CreatePage(jsonPath+"/createPage.json", args[0],
			api.PROPERTIES_SELECT, "table_name", args[1],
			api.PROPERTIES_TITLE, "column_name", args[2],
			api.PROPERTIES_NUMBER, "ordinal_position", v3,
			api.PROPERTIES_TEXT, "column_default", args[4],
			api.PROPERTIES_SELECT, "is_nullable", args[5],
			api.PROPERTIES_TEXT, "data_type", args[6],
			api.PROPERTIES_NUMBER, "character_maximum_length", v7)
		if err != nil {
			return err
		}
		log.Printf("Page %s inserted\n", pg)
	} else {
		log.Fatal("insufficient arguments provided")
	}

	return nil
}
