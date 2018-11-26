package teamtailor-integration-go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Candidate struct {
	Email       string    `json:"email" bson:"email"`
	Connected   bool      `json:"connected" bson:"connected"`
	Created     time.Time `json:"created-at" bson:"created-at"` //TODO: Should be date format in json
	Firstname   string    `json:"first-name" bson:"first-name"`
	Lastname    string    `json:"last-name" bson:"last-name"`
	LinkedinUID string    `json:"linkedin-uid" bson:"linkedin-uid"`
	FacebookUID string    `json:"facebook-id" bson:"facebook-id"`
	Phone       string    `json:"phone" bson:"phone"`
	Picture     string    `json:"picture" bson:"picture"`
	Pitch       string    `json:"pitch" bson:"pitch"`
	Sourced     bool      `json:"sourced" bson:"sourced"`
	Tags        []string  `json:"tags" bson:"tags"`
	UpdatedAt   time.Time `json:"updated-at" bson:"updated-at"`
}

// Convert Candidate struct into JSON
func candidateToJSON(cand Candidate) []byte {
	data, err := json.Marshal(cand)
	if err != nil {
		return nil
	}
	return data
}

// PostCandidate creates and executes a POST-request to the TeamTailor API and returns the resposne body as a []byte
func (t *TeamTailor) PostCandidate(c Candidate) ([]byte, error) {
	client := &http.Client{}

	// TODO: Check token validity

	cand := candidateToJSON(c)
	postData := bytes.NewReader(cand)

	req, _ := http.NewRequest("POST", t.ApiHost+"candidate", postData)
	req.Header.Set("Authorization", "Token token="+t.Authorization)
	req.Header.Set("X-Api-Version", t.APIversion)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	return body, nil

}

// func GetCandidate
func GetCandidate() (string, error) {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/users/1")

	if err != nil {
		return "fail", nil
	}

	cand := Candidate{}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &cand)

	cand.Firstname = "Andreas"
	cand.Lastname = "Mann"

	fmt.Println(cand)

	return "success", nil
}

// func GetCandidates

// func DeleteCandidate
