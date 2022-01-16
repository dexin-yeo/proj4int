package api

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"notion/client"
)

type Fields int

const (
	PROPERTIES_NUMBER Fields = iota
	PROPERTIES_TEXT
	PROPERTIES_TITLE
	PROPERTIES_SELECT
	PARENT_DATABASE
	PARENT_PAGE
	FIELDS_TITLE
)

type API struct {
	client *client.Client
}

func CreateAPI() *API {
	notionKey := os.Getenv("NOTION_KEY")
	auth := "Bearer " + notionKey

	header := http.Header{
		"Authorization":  []string{auth},
		"Notion-version": []string{"2021-08-16"},
	}
	client := client.CreateClient("https://api.notion.com/v1/", header)

	return &API{client}
}

func (api *API) ListDatabases() {
	ret, err := api.client.SendGet("databases")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(ret)
	// type Results struct {
	// 	Object      string        `json:object`
	// 	Id          string        `json:id`
	// 	Cover       string        `json:cover`
	// 	Icon        string        `json:icon`
	// 	CreatedTime string        `json:created_time`
	// 	LastEdit    string        `json:last_edited_time`
	// 	Title       []interface{} `json:title`
	// 	Propertiles interface{}   `json:properties`
	// 	Parent      interface{}   `json:parent`
	// 	Url         string        `json:url`
	// }

	type Return struct {
		Object string `json:object`
		// Result []Results `json:results`
		Results []interface{} `json:results`
		// Results json.RawMessage
		Next string `json:next_cursor`
		More bool   `json:has_more`
	}

	var res Return
	err = json.Unmarshal([]byte(ret), &res)
	if err != nil {
		log.Fatal(err)
	}

	list := make([]string, len(res.Results))

	for i, v := range res.Results {
		list[i] = v.(map[string]interface{})["id"].(string)
	}

	fmt.Println(list)
}

func (api *API) post(file, method string, parent Fields, id string, params ...interface{}) (string, error) {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	var rawMap map[string]interface{}
	// unmarshal into unstructured map
	err = json.Unmarshal(input, &rawMap)
	if err != nil {
		return "", err
	}

	// edit parent
	switch parent {
	case PARENT_DATABASE:
		rawMap = editParent(rawMap, id, PARENT_DATABASE)
	case PARENT_PAGE:
		rawMap = editParent(rawMap, id, PARENT_PAGE)
	}

	// at least 3 additional condition given or 2 + 1(empty)
	if len(params) >= 3 {
		// edit any necessary fields
		for i := 0; i < len(params); i += 3 {
			switch params[i].(Fields) {
			case FIELDS_TITLE:
				rawMap = editTitle(rawMap, params[i+1].(string))
			case PROPERTIES_NUMBER:
				rawMap = editNumberProperties(rawMap, params[i+1].(string), params[i+2].(int))
			case PROPERTIES_TEXT:
				rawMap = editTextTitleProperties(rawMap, params[i+1].(string), params[i+2].(string), PROPERTIES_TEXT)
			case PROPERTIES_TITLE:
				rawMap = editTextTitleProperties(rawMap, params[i+1].(string), params[i+2].(string), PROPERTIES_TITLE)
			case PROPERTIES_SELECT:
				rawMap = editSelectProperties(rawMap, params[i+1].(string), params[i+2].(string))
			}
		}
	}

	// marshal back to json string
	input, err = json.Marshal(rawMap)
	if err != nil {
		return "", err
	}

	// send the req
	ret, err := api.client.Send("POST", method, input)
	if err != nil {
		return "", err
	}

	type Return struct {
		Id string `json:id`
	}

	var res Return
	// unmarshal to get the page/database ID created/edited
	err = json.Unmarshal([]byte(ret), &res)
	if err != nil {
		return "", err
	}

	return res.Id, nil
}

var (
	numProp       string = "\"%s\":{\"number\":%d}"
	textTitleProp string = "\"%s\":{\"%s\":[{\"text\":{\"content\":\"%s\"}}]}"
	selectProp    string = "\"%s\":{\"select\":{\"name\":\"%s\"}}"
)

func (api *API) patch(endpoint string, params ...interface{}) (string, error) {
	// rough format of the update json for notion
	data := "{\"properties\":{%s}}"
	props := ""

	// at least 3 additional condition
	if len(params) >= 3 {
		// get the last element index
		lastEle := len(params) - 3

		// edit any necessary fields
		for i := 0; i < len(params); i += 3 {
			switch params[i].(Fields) {
			case PROPERTIES_NUMBER:
				props += fmt.Sprintf(numProp, params[i+1].(string), params[i+2].(int))
			case PROPERTIES_TEXT:
				props += fmt.Sprintf(textTitleProp, params[i+1].(string), "rich_text", params[i+2].(string))
			case PROPERTIES_TITLE:
				props += fmt.Sprintf(textTitleProp, params[i+1].(string), "title", params[i+2].(string))
			case PROPERTIES_SELECT:
				props += fmt.Sprintf(selectProp, params[i+1].(string), params[i+2].(string))
			}

			// adds ',' into the json string if not last property
			if i != lastEle {
				props += ","
			}
		}
	} else {
		return "", fmt.Errorf("not enough argument")
	}

	data = fmt.Sprintf(data, props)

	// send the req
	ret, err := api.client.Send(http.MethodPatch, endpoint, []byte(data))
	if err != nil {
		return "", err
	}

	type Return struct {
		Id string `json:id`
	}

	var res Return
	// unmarshal to get the page/database ID created/edited
	err = json.Unmarshal([]byte(ret), &res)
	if err != nil {
		return "", err
	}

	return res.Id, nil
}

func (api *API) CreateDatabase(file string, pageID string, title string) (string, error) {
	return api.post(file, "databases", PARENT_PAGE, pageID, FIELDS_TITLE, title, "")
}

func (api *API) CreatePage(file string, databaseID string, params ...interface{}) (string, error) {
	return api.post(file, "pages", PARENT_DATABASE, databaseID, params...)
}

func (api *API) UpdatePage(pageID string, params ...interface{}) (string, error) {
	return api.patch("pages/"+pageID, params...)
}

func editParent(src map[string]interface{}, input string, pType Fields) map[string]interface{} {
	dType := ""
	switch pType {
	case PARENT_DATABASE:
		dType = "database_id"
	case PARENT_PAGE:
		dType = "page_id"
	}
	src["parent"].(map[string]interface{})[dType] = input
	return src
}

func editTitle(src map[string]interface{}, input string) map[string]interface{} {
	src["title"].([]interface{})[0].(map[string]interface{})["text"].(map[string]interface{})["content"] = input
	return src
}

func editNumberProperties(src map[string]interface{}, colname string, input int) map[string]interface{} {
	src["properties"].(map[string]interface{})[colname].(map[string]interface{})["number"] = input
	return src
}

func editTextTitleProperties(src map[string]interface{}, colname string, input string, propType Fields) map[string]interface{} {
	dType := ""
	switch propType {
	case PROPERTIES_TEXT:
		dType = "rich_text"
	case PROPERTIES_TITLE:
		dType = "title"
	}
	src["properties"].(map[string]interface{})[colname].(map[string]interface{})[dType].([]interface{})[0].(map[string]interface{})["text"].(map[string]interface{})["content"] = input
	return src
}

func editSelectProperties(src map[string]interface{}, colname string, input string) map[string]interface{} {
	src["properties"].(map[string]interface{})[colname].(map[string]interface{})["select"].(map[string]interface{})["name"] = input
	return src
}

/*
Reads the csv file and extract the data
Param:
	- file: file path to the csv file
Returns:
	- array of rows in string
*/
func ReadCSVFile(file string) [][]string {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return records
}
