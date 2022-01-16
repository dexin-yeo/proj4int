package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	endpoint   string
	getHeader  http.Header
	postHeader http.Header
}

func CreateClient(endpoint string, header http.Header) *Client {
	pHeader := header.Clone()
	pHeader.Add("Content-Type", "application/json")
	return &Client{endpoint, header, pHeader}
}

func (client *Client) SendGet(object string) (string, error) {
	url := client.endpoint + object
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header = client.getHeader

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("request failed with status %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (client *Client) Send(method, object string, msg []byte) (string, error) {
	url := client.endpoint + object
	req, err := http.NewRequest(method, url, bytes.NewBuffer(msg))
	if err != nil {
		return "", err
	}
	req.Header = client.postHeader

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	// if res.StatusCode != 200 {
	// 	return "", fmt.Errorf("request failed with status %d", res.StatusCode)
	// }
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// func main() {
// 	url := "https://api.notion.com/v1/"
// 	dbUrl := url + "databases"
// 	fmt.Println(dbUrl)
// 	req, err := http.NewRequest("GET", dbUrl, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	notionKey := os.Getenv("NOTION_KEY")
// 	auth := "Bearer " + notionKey
// 	fmt.Println(auth)

// 	req.Header.Add("Authorization", auth)
// 	req.Header.Add("Notion-version", "2021-08-16")
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(string(body))
// }
