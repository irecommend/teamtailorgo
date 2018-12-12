package teamtailorgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/manyminds/api2go/jsonapi"
)

type Candidate struct {
	Email       string    `json:"email" jsonapi:"email"`
	Connected   bool      `json:"connected" jsonapi:"connected"`
	Created     time.Time `json:"created-at" jsonapi:"created-at"` //TODO: Should be date format in json
	Firstname   string    `json:"first-name" jsonapi:"first-name"`
	Lastname    string    `json:"last-name" jsonapi:"last-name"`
	LinkedinUID string    `json:"linkedin-uid" jsonapi:"linkedin-uid"`
	LinkedinURL string    `json:"linkedin-url" jsonapi:"linkedin-url"`
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

type CandidateResponseData struct {
	ID         string               `json:"id"`
	Type       string               `json:"type"`
	Links      *Links               `json:"links"`
	Attributes *CandidateAttributes `json:"attributes"`
}

type Links struct {
	Self string `json:"self"`
}

type CandidateAttributes struct {
	Email           string    `json:"email" jsonapi:"email"`
	Connected       bool      `json:"connected" jsonapi:"connected"`
	Created         time.Time `json:"created-at" jsonapi:"created-at"` //TODO: Should be date format in json
	Firstname       string    `json:"first-name" jsonapi:"first-name"`
	Lastname        string    `json:"last-name" jsonapi:"last-name"`
	LinkedinUID     string    `json:"linkedin-uid" jsonapi:"linkedin-uid"`
	LinkedinURL     string    `json:"linkedin-url" jsonapi:"linkedin-url"`
	FacebookUID     string    `json:"facebook-id" jsonapi:"facebook-id"`
	Phone           string    `json:"phone" jsonapi:"phone"`
	Picture         string    `json:"picture" jsonapi:"picture"`
	Pitch           string    `json:"pitch" jsonapi:"pitch"`
	Sourced         bool      `json:"sourced" jsonapi:"sourced"`
	Tags            []string  `json:"tags" jsonapi:"tags"`
	UpdatedAt       time.Time `json:"updated-at" jsonapi:"updated-at"`
	ReferringSite   string    `json:"referring-site" jsonapi:"referringsite"`
	ReferringURL    string    `json:"referring-url" jsonapi:"referring-url"`
	Resume          string    `json:"resume" jsonapi:"resume"`
	Unsubscribed    bool      `json:"unsubscribed" jsonapi:"unsubscribed"`
	FacebookProfile string    `json:"facebook-profile" jsonapi:"facebook-profile"`
	LinkedinProfile string    `json:"linkedin-profile" jsonapi:"linkedin-profile"`
}

// TODO: RETURN EVERYTHING FROM PACKAGE
type CandidateResponse struct {
	Data CandidateResponseData `json:"data"`
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
func (t *TeamTailor) PostCandidate(c Candidate) (CandidateResponse, error) {

	cand := candidateToJSON(c)
	postData := bytes.NewReader(cand)
	var rc CandidateResponse

	req, _ := http.NewRequest("POST", baseUrl+"candidates", postData)
	req.Header.Set("Authorization", "Token token="+t.Token)
	req.Header.Set("X-Api-Version", apiVersion)
	req.Header.Set("Content-Type", contentType)

	resp, err := t.HTTPClient.Do(req)
	if err != nil {
		return rc, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	// TODO: ERROR HANDLING

	err = json.Unmarshal(body, &rc)

	if err != nil {
		return rc, err
	}

	defer resp.Body.Close()

	return rc, nil

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
