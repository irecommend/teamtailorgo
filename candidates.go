package teamtailorgo

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"time"

	japi "github.com/google/jsonapi"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/pkg/errors"
)

type CandidateRequest struct {
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
	Type      string            `json:"type"`
	Candidate *CandidateRequest `json:"attributes"`
}

type Candidate struct {
	ID              string   `json:"-" jsonapi:"primary,candidates"`
	Email           string   `json:"email" jsonapi:"attr,email"`
	Connected       bool     `json:"connected" jsonapi:"attr,connected"`
	Created         string   `json:"created-at" jsonapi:"attr,created-at"`
	Firstname       string   `json:"first-name" jsonapi:"attr,first-name"`
	Lastname        string   `json:"last-name" jsonapi:"attr,last-name"`
	LinkedinUID     string   `json:"linkedin-uid" jsonapi:"attr,linkedin-uid"`
	LinkedinURL     string   `json:"linkedin-url" jsonapi:"attr,linkedin-url"`
	FacebookUID     string   `json:"facebook-id" jsonapi:"attr,facebook-id"`
	Phone           string   `json:"phone" jsonapi:"attr,phone"`
	Picture         string   `json:"picture" jsonapi:"attr,picture"`
	Pitch           string   `json:"pitch" jsonapi:"attr,pitch"`
	Sourced         bool     `json:"sourced" jsonapi:"attr,sourced"`
	Tags            []string `json:"tags" jsonapi:"attr,tags"`
	UpdatedAt       string   `json:"updated-at" jsonapi:"attr,updated-at"`
	ReferringSite   string   `json:"referring-site" jsonapi:"attr,referringsite"`
	ReferringURL    string   `json:"referring-url" jsonapi:"attr,referring-url"`
	Resume          string   `json:"resume" jsonapi:"attr,resume"`
	Unsubscribed    bool     `json:"unsubscribed" jsonapi:"attr,unsubscribed"`
	FacebookProfile string   `json:"facebook-profile" jsonapi:"attr,facebook-profile"`
	LinkedinProfile string   `json:"linkedin-profile" jsonapi:"attr,linkedin-profile"`
}

// Convert Candidate struct into JSON
func candidateToJSON(cand CandidateRequest) ([]byte, error) {

	// Use external package that sadly forces and ID on the JSON object
	data, err := jsonapi.Marshal(cand)
	if err != nil {
		return nil, err
	}

	// Unmarshal back to custom struct to remove ID
	unmrsh := CandidateJSONApi{}
	err = json.Unmarshal(data, &unmrsh)
	if err != nil {
		return nil, err
	}
	unmrsh.Data.Type = "candidates"

	// Marshal custom struct
	postData, err := json.Marshal(unmrsh)
	if err != nil {
		return nil, err
	}

	return postData, nil
}

// PostCandidate creates and executes a POST-request to the TeamTailor API and returns the resposne body as a []byte
// TODO: Should return existing candidate if that is the case
func (t *TeamTailor) PostCandidate(c CandidateRequest) (*Candidate, error) {

	var rc Candidate

	cand, err := candidateToJSON(c)
	if err != nil {
		return &rc, errors.New("Invalid structure of provided candidate")
	}

	postData := bytes.NewReader(cand)

	req, _ := http.NewRequest("POST", baseURL+"candidates", postData)
	req.Header.Set("Authorization", "Token token="+t.Token)
	req.Header.Set("X-Api-Version", apiVersion)
	req.Header.Set("Content-Type", contentType)

	resp, err := t.HTTPClient.Do(req)
	if err != nil {
		return &rc, err
	}

	// New candidate posted
	if resp.StatusCode == 201 {
		err = japi.UnmarshalPayload(resp.Body, &rc)
		if err != nil {
			log.Println("ERROR IN UNMARSHAL", err)
			return &rc, err
		}

		defer resp.Body.Close()

		return &rc, nil
	} else if resp.StatusCode == 422 {
		// Candidate existed in TeamTailor
		cand, err := t.GetCandidateByEmail(c.Email)
		if err != nil {
			return &rc, err
		}

		// TODO: Update candidate with irec tag

		return cand, nil
	} else {
		return &rc, errors.New("Failed posting candidate")
	}
}

// func GetCandidate
func (t *TeamTailor) GetCandidate(id string) (Candidate, error) {

	var cand Candidate
	req, _ := http.NewRequest("GET", baseURL+"candidates/"+id, nil)
	t.SetHeaders(req)

	resp, err := t.HTTPClient.Do(req)
	if err != nil {
		return cand, err
	}

	err = japi.UnmarshalPayload(resp.Body, &cand)
	if err != nil {
		return cand, err
	}

	defer resp.Body.Close()

	return cand, nil
}

// func GetCandidateByEmail
func (t *TeamTailor) GetCandidateByEmail(email string) (*Candidate, error) {

	var candidate *Candidate
	// Get candidates
	candidates, err := t.GetCandidates()
	if err != nil {
		return candidate, err
	}

	// Filter through candidates by email
	for _, cand := range candidates {
		if cand.Email == email {
			candidate = cand
			break
		}
	}

	if candidate.Email == "" {
		return candidate, errors.New("Could not find existing candidate by email")
	}

	// Return the one that matches
	return candidate, nil

}

// func GetCandidates
func (t *TeamTailor) GetCandidates() ([]*Candidate, error) {
	var cands []*Candidate

	req, _ := http.NewRequest("GET", baseURL+"candidates", nil)
	req.Header.Set("Authorization", "Token token="+t.Token)
	req.Header.Set("X-Api-Version", apiVersion)
	req.Header.Set("Content-Type", contentType)

	resp, err := t.HTTPClient.Do(req)
	if err != nil {
		return cands, err
	}

	candidates, err := japi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Candidate)))
	if err != nil {
		return cands, err
	}

	for _, candidate := range candidates {
		c, _ := candidate.(*Candidate)
		cands = append(cands, c)
	}

	defer resp.Body.Close()

	return cands, nil
}

// func (t *TeamTailor) UpdateCandidate(id string) error {

// }

// func DeleteCandidate

// func GetCandidateJobApplications

// func CreateCandidateJobApplication

// func SetHeaders
func (t *TeamTailor) SetHeaders(r *http.Request) {
	r.Header.Set("Authorization", "Token token="+t.Token)
	r.Header.Set("X-Api-Version", apiVersion)
	r.Header.Set("Content-Type", contentType)
}

// JSON API INTERFACE FUNCTIONS
func (c *Candidate) SetToOneReferenceID(name, ID string) error {
	c.ID = ID
	return nil
}

func (c CandidateRequest) GetID() string {
	return "0"
}

func (c *Candidate) SetID(ID string) error {
	c.ID = ID
	return nil
}

func (c Candidate) GetID() string {
	return c.ID
}
