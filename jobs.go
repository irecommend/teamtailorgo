package teamtailorgo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/manyminds/api2go/jsonapi"
)

// MetaTags checks how many job objects and how many pages exist in the
// Teamtailor instance. Used if more than 30 jobs exist for a company.
type MetaTags struct {
	RecordCount int `json:"record-count"`
	PageCount   int `json:"page-count"`
}

// Meta is a wrapper for meta tags. Needed in order to Marshal()
type Meta struct {
	MetaTags *MetaTags `json:"meta"`
}

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

// GetFirstJobPage fetches the 30 jobs that are on the first page of Teamtailor response.
// If there are more than 30 jobs to fetch from Teamtailor, use GetAllJobs().
func (t TeamTailor) GetFirstJobPage() ([]Job, error) {

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

// GetAllJobs returns all jobs for a company. Used if page count in Teamtailor is greater
// than one and multiple GET-requests are needed to get all jobs.
func (t TeamTailor) GetAllJobs() ([]Job, error) {

	var jobs []Job
	var meta Meta

	// Begin request to get Meta tags to find record and page count.
	req, _ := http.NewRequest("GET", baseURL+"jobs", nil)
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

	err = json.Unmarshal(body, &meta)
	if err != nil {
		return jobs, err
	}

	// Begin request to get all pages of jobs
	for n := 1; n <= meta.MetaTags.PageCount; n++ {
		req, _ := http.NewRequest("GET", baseURL+"jobs?page%5Bsize%5D=30&page%5Bnumber%5D="+strconv.Itoa(n), nil)
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
