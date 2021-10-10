package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	address := flag.String("server", "http://localhost:8080", "HTTP gateway url, e.g. http://localhost:8080")
	flag.Parse()

	t := time.Now().In(time.UTC)
	pfx := t.Format(time.RFC3339Nano)

	// Create
	resp, err := http.Post(*address+"/api/v1/foo", "application/json", strings.NewReader(fmt.Sprintf(`
		{
			"api_version":"v1",
			"foo": {
				"title":"title (%s)",
				"desc": "desc (%s)"
			}
		}
	`, pfx, pfx)))
	if err != nil {
		log.Fatalf("[ERROR] Failed to call Create method: %v", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var body string
	if err != nil {
		body = fmt.Sprintf("[ERROR] Failed to read Create response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("[INFO] Create response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	var created struct {
		API_VERSION string `json:"api_version"`
		ID          string `json:"id"`
	}

	err = json.Unmarshal(bodyBytes, &created)
	if err != nil {
		log.Fatalf("[ERROR] Failed to unmarshal JSON response of Create method: %v", err)
		fmt.Println("error:", err)
	}

	// Read
	resp, err = http.Get(fmt.Sprintf("%s%s/%s", *address, "/api/v1/foo", created.ID))
	if err != nil {
		log.Fatalf("[ERROR] Failed to call Read method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("[ERROR] Failed read Read response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Read response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// Update
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s", *address, "/api/v1/foo", created.ID),
		strings.NewReader(fmt.Sprintf(`
		{
			"api_version":"v1",
			"foo": {
				"title":"title (%s) + updated",
				"desc":"desc (%s) + updated"
			}
		}
	`, pfx, pfx)))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("[ERROR] Failed to call Update method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("[ERROR] Failed read Update response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("[INFO] Update response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// Delete
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("%s%s/%s", *address, "/api/v1/foo", created.ID), nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("[ERROR] Failed to call Delete method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("[ERROR] Failed read Delete response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("[INFO] Delete response: Code=%d, Body=%s\n\n", resp.StatusCode, body)
}
