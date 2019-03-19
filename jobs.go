package teamtailorgo

import (
	"io/ioutil"
	"net/http"

	"github.com/manyminds/api2go/jsonapi"
)

// TODO: Add picture which consists of substrings "standard" and "thumb" or "" if no picture for job
type Job struct {
	ID          string
	Type        string   `json:"type"`
	Body        string   `json:"body"`
	EndDate     string   `json:"end-date"`
	HumanStatus string   `json:"human-status"`
	Internal    bool     `json:"internal"`
	Pinned      bool     `json:"pinned"`
	StartDate   string   `json:"start-date"`
	Status      string   `json:"status"`
	Title       string   `json:"title"`
	Tags        []string `json:"tags"`
	Created     string   `json:"created-at"`
}

func (t TeamTailor) GetAllJobs() ([]Job, error) {

	var jobs []Job

	req, _ := http.NewRequest("GET", baseURL+"jobs?page%5Bsize%5D=30", nil)
	req.Header.Set("Authorization", "Token token="+t.Token)
	req.Header.Set("X-Api-Version", apiVersion)
	req.Header.Set("Content-Type", contentType)

	resp, err := t.HTTPClient.Do(req)
	if err != nil {
		return jobs, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return jobs, err
	}

	err = jsonapi.Unmarshal(body, &jobs)
	if err != nil {
		return jobs, err
	}

	return jobs, nil
}

// GetJob
func (t TeamTailor) GetJob(id string) (Job, error) {

	var job Job

	req, _ := http.NewRequest("GET", baseURL+"jobs/"+id, nil)
	req.Header.Set("Authorization", "Token token="+t.Token)
	req.Header.Set("X-Api-Version", apiVersion)
	req.Header.Set("Content-Type", contentType)

	resp, err := t.HTTPClient.Do(req)
	if err != nil {
		return job, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return job, err
	}

	err = jsonapi.Unmarshal(body, &job)
	if err != nil {
		return job, err
	}

	return job, nil
}

// JSONAPI functions

func (j *Job) SetID(ID string) error {
	j.ID = ID
	return nil
}

func (j Job) GetID() string {
	return j.ID
}

func (j Job) SetToOneReferenceID(name, ID string) error {
	j.ID = ID
	return nil
}
