package teamtailorgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/manyminds/api2go/jsonapi"
)

type MarshalIdentifier interface {
	GetID() string
}

type Candidate struct {
	Email       string    `json:"email" jsonapi:"email"`
	Connected   bool      `json:"connected" jsonapi:"connected"`
	Created     time.Time `json:"created-at" jsonapi:"created-at"` //TODO: Should be date format in json
	Firstname   string    `json:"first-name" jsonapi:"first-name"`
	Lastname    string    `json:"last-name" jsonapi:"last-name"`
	LinkedinUID string    `json:"linkedin-uid" jsonapi:"linkedin-uid"`
	FacebookUID string    `json:"facebook-id" jsonapi:"facebook-id"`
	Phone       string    `json:"phone" jsonapi:"phone"`
	Picture     string    `json:"picture" jsonapi:"picture"`
	Pitch       string    `json:"pitch" jsonapi:"pitch"`
	Sourced     bool      `json:"sourced" jsonapi:"sourced"`
	Tags        []string  `json:"tags" jsonapi:"tags"`
	UpdatedAt   time.Time `json:"updated-at" jsonapi:"updated-at"`
}

type CandidateJSONApi struct {
	Data *CandidateConverted `json:"data"`
}

type CandidateConverted struct {
	Type      string     `json:"type"`
	Candidate *Candidate `json:"attributes"`
}

func (c Candidate) GetID() string {
	return "0"
}

// Convert Candidate struct into JSON
func candidateToJSON(cand Candidate) []byte {

	// Use external package that sadly forces and ID on the JSON object
	data, err := jsonapi.Marshal(cand)
	if err != nil {
		return nil
	}

	// Unmarshal back to custom struct to remove ID
	unmrsh := CandidateJSONApi{}
	err = json.Unmarshal(data, &unmrsh)
	if err != nil {
		return nil
	}

	// Marshal custom struct
	postData, err := json.Marshal(unmrsh)
	if err != nil {
		return nil
	}

	return postData
}

// PostCandidate creates and executes a POST-request to the TeamTailor API and returns the resposne body as a []byte
func (t *TeamTailor) PostCandidate(c Candidate) ([]byte, error) {
	client := &http.Client{}

	// TODO: Check token validity

	cand := candidateToJSON(c)
	fmt.Println(string(cand))
	postData := bytes.NewReader(cand)

	req, _ := http.NewRequest("POST", t.ApiHost+"candidates", postData)
	req.Header.Set("Authorization", "Token token="+t.Authorization)
	req.Header.Set("X-Api-Version", t.APIversion)
	req.Header.Set("Content-Type", "application/vnd.api+json")

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
