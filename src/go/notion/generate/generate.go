package generate

import (
	"log"
	"notion/api"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
Create and populate table-definition.csv into notion
Param:
	- args: args from command line
Returns:
	- error
*/
func Generate(args []string) error {
	// checks whether page id is given
	if len(args) < 1 {
		log.Fatal("page id is not given")
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

	// reads and extract data from the csv the csv file
	data := api.ReadCSVFile(jsonPath + "/table-definition.csv")

	obj := api.CreateAPI()
	log.Printf("Creating database to page %s\n", args[0])
	id, err := obj.CreateDatabase(jsonPath+"/createDatabase.json", args[0], "table-definition")
	if err != nil {
		return err
	}
	log.Printf("Database %s created to page %s\n", id, args[0])

	// extra symbols
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		return err
	}

	// first element contains the col names
	cols := data[0]
	// remove extra symbols from title name to prevent errs
	for i, v := range cols {
		cols[i] = re.ReplaceAllString(v, "")
	}

	log.Printf("Populating rows to database %s as pages...\n", id)
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
		pg, err := obj.CreatePage(jsonPath+"/createPage.json", id,
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

	return nil
}
