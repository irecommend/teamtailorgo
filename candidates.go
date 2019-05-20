package teamtailorgo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"time"

	japi "github.com/google/jsonapi"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/pkg/errors"
)

type CandidateRequest struct {
	Email       string    `json:"email" jsonapi:"attr, email"`
	Connected   bool      `json:"connected" jsonapi:"attr,connected"`
	Created     time.Time `json:"created-at" jsonapi:"attr, created-at"`
	Firstname   string    `json:"first-name" jsonapi:"attr, first-name"`
	Lastname    string    `json:"last-name" jsonapi:"attr, last-name"`
	LinkedinUID string    `json:"linkedin-uid" jsonapi:"attr, linkedin-uid"`
	LinkedinURL string    `json:"linkedin-url" jsonapi:"attr, linkedin-url"`
	FacebookUID string    `json:"facebook-id" jsonapi:"attr, facebook-id"`
	Phone       string    `json:"phone" jsonapi:"attr, phone"`
	Picture     string    `json:"picture" jsonapi:"attr, picture"`
	Pitch       string    `json:"pitch" jsonapi:"attr, pitch"`
	Sourced     bool      `json:"sourced" jsonapi:"attr, sourced"`
	Tags        []string  `json:"tags" jsonapi:"attr, tags"`
	UpdatedAt   time.Time `json:"updated-at" jsonapi:"attr, updated-at"`
}

type CandidateJSONApi struct {
	Data *CandidateConverted `json:"data"`
}

type CandidateConverted struct {
	Type      string            `json:"type"`
	Candidate *CandidateRequest `json:"attributes"`
}

type CandidateRequestResume struct {
	Email       string    `json:"email" jsonapi:"email"`
	Connected   bool      `json:"connected" jsonapi:"connected"`
	Created     time.Time `json:"created-at" jsonapi:"created-at"`
	Firstname   string    `json:"first-name" jsonapi:"first-name"`
	Lastname    string    `json:"last-name" jsonapi:"last-name"`
	LinkedinUID string    `json:"linkedin-uid" jsonapi:"linkedin-uid"`
	LinkedinURL string    `json:"linkedin-url" jsonapi:"linkedin-url"`
	FacebookUID string    `json:"facebook-id" jsonapi:"facebook-id"`
	Phone       string    `json:"phone" jsonapi:"phone"`
	Picture     string    `json:"picture" jsonapi:"picture"`
	Pitch       string    `json:"pitch" jsonapi:"pitch"`
	Resume      string    `json:"resume" jsonapi:"resume"`
	Sourced     bool      `json:"sourced" jsonapi:"sourced"`
	Tags        []string  `json:"tags" jsonapi:"tags"`
	UpdatedAt   time.Time `json:"updated-at" jsonapi:"updated-at"`
}

type CandidateResumeJSONApi struct {
	Data *CandidateResumeConverted `json:"data"`
}

type CandidateResumeConverted struct {
	Type      string                  `json:"type"`
	Candidate *CandidateRequestResume `json:"attributes"`
}

type Candidate struct {
	ID           string   `json:"-" jsonapi:"primary,candidates"`
	Email        string   `json:"email" jsonapi:"attr,email"`
	Connected    bool     `json:"connected" jsonapi:"attr,connected"`
	Created      string   `json:"created-at" jsonapi:"attr,created-at"`
	Firstname    string   `json:"first-name" jsonapi:"attr,first-name"`
	Lastname     string   `json:"last-name" jsonapi:"attr,last-name"`
	LinkedinUID  string   `json:"linkedin-uid" jsonapi:"attr,linkedin-uid"`
	LinkedinURL  string   `json:"linkedin-url" jsonapi:"attr,linkedin-url"`
	FacebookUID  string   `json:"facebook-id" jsonapi:"attr,facebook-id"`
	Phone        string   `json:"phone" jsonapi:"attr,phone"`
	Picture      string   `json:"picture" jsonapi:"attr,picture"`
	Pitch        string   `json:"pitch" jsonapi:"attr,pitch"`
	Sourced      bool     `json:"sourced" jsonapi:"attr,sourced"`
	Tags         []string `json:"tags" jsonapi:"attr,tags"`
	UpdatedAt    string   `json:"updated-at" jsonapi:"attr,updated-at"`
	ReferringURL string   `json:"referring-url" jsonapi:"attr,referring-url"`
	Resume       string   `json:"resume" jsonapi:"attr,resume"`
	Unsubscribed bool     `json:"unsubscribed" jsonapi:"attr,unsubscribed"`
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

// Convert Candidate struct into JSON
func candidateResumeToJSON(cand CandidateRequestResume) ([]byte, error) {

	// Use external package that sadly forces and ID on the JSON object
	data, err := jsonapi.Marshal(cand)
	if err != nil {
		return nil, err
	}

	// Unmarshal back to custom struct to remove ID
	unmrsh := CandidateResumeJSONApi{}
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
	defer resp.Body.Close()

	// New candidate posted
	if resp.StatusCode == 201 {
		err = japi.UnmarshalPayload(resp.Body, &rc)
		if err != nil {
			return &rc, err
		}

		return &rc, nil
	} else if resp.StatusCode == 422 {
		// Candidate existed in TeamTailor
		cand, err := t.GetCandidateByEmail(c.Email)
		if err != nil {
			return &rc, err
		}

		return cand, nil
	} else {
		return &rc, errors.New("Failed posting candidate, got response code " + strconv.Itoa(resp.StatusCode))
	}
}

// PostCandidateResume executes POST-request with a signedURL to access a candidate's resume
func (t *TeamTailor) PostCandidateResume(c CandidateRequestResume) (*Candidate, error) {

	var rc Candidate

	cand, err := candidateResumeToJSON(c)
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
	defer resp.Body.Close()

	// New candidate posted
	if resp.StatusCode == 201 {
		err = japi.UnmarshalPayload(resp.Body, &rc)
		if err != nil {
			return &rc, err
		}

		return &rc, nil
	} else if resp.StatusCode == 422 {
		// Candidate existed in TeamTailor
		cand, err := t.GetCandidateByEmail(c.Email)
		if err != nil {
			return &rc, err
		}

		return cand, nil
	} else {
		return &rc, errors.New("Failed posting candidate, got response code " + strconv.Itoa(resp.StatusCode))
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
	defer resp.Body.Close()

	err = japi.UnmarshalPayload(resp.Body, &cand)
	if err != nil {
		return cand, err
	}

	return cand, nil
}

// func GetCandidateByEmail
func (t *TeamTailor) GetCandidateByEmail(email string) (*Candidate, error) {

	var cand *Candidate
	req, _ := http.NewRequest("GET", baseURL+"candidates?filter[email]="+email, nil)
	t.SetHeaders(req)

	resp, err := t.HTTPClient.Do(req)
	if err != nil {
		return cand, err
	}
	defer resp.Body.Close()

	candidates, err := japi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Candidate)))
	if err != nil {
		return cand, err
	}

	cand = candidates[0].(*Candidate)

	return cand, nil

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
	defer resp.Body.Close()

	candidates, err := japi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Candidate)))
	if err != nil {
		return cands, err
	}

	for _, candidate := range candidates {
		c, _ := candidate.(*Candidate)
		cands = append(cands, c)
	}

	return cands, nil
}

func (t *TeamTailor) UpdateCandidate(c Candidate) error {

	data, err := jsonapi.Marshal(c)
	if err != nil {
		return err
	}
	postData := bytes.NewReader(data)

	req, _ := http.NewRequest("PATCH", baseURL+"candidates/"+c.ID, postData)
	t.SetHeaders(req)

	resp, err := t.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		return nil
	}

	return errors.New("Failed updating candidate, got response code " + strconv.Itoa(resp.StatusCode))
}

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

func (c CandidateRequestResume) GetID() string {
	return "0"
}

func (c *Candidate) SetID(ID string) error {
	c.ID = ID
	return nil
}

func (c Candidate) GetID() string {
	return c.ID
}
