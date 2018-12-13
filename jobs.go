package teamtailorgo

import (
	"io/ioutil"
	"net/http"

	"github.com/manyminds/api2go/jsonapi"
)

type Job struct {
	ID          string
	Type        string   `json:"type"`
	Body        string   `json:"body"`
	EndDate     string   `json:"end-date"`
	HumanStatus string   `json:"human-status"`
	Internal    bool     `json:"internal"`
	Picture     string   `json:"picture"`
	Pinned      bool     `json:"pinned"`
	StartDate   string   `json:"start-date"`
	Status      string   `json:"status"`
	Title       string   `json:"title"`
	Tags        []string `json:"tags"`
	Created     string   `json:"created-at"`
}

func (t TeamTailor) GetAllJobs() ([]Job, error) {

	var jobs []Job

	req, _ := http.NewRequest("GET", baseURL+"jobs", nil)
	req.Header.Set("Authorization", "Token token="+t.Token)
	req.Header.Set("X-Api-Version", apiVersion)
	req.Header.Set("Content-Type", contentType)

	resp, err := t.HTTPClient.Do(req)
	if err != nil {
		return jobs, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	// TODO: ERROR HANDLING

	err = jsonapi.Unmarshal(body, &jobs)
	if err != nil {
		return jobs, err
	}

	defer resp.Body.Close()

	return jobs, nil
}

func (j *Job) SetID(ID string) error {
	j.ID = ID
	return nil
}

func (j Job) GetID() string {
	return string(j.ID)
}

func (j Job) SetToOneReferenceID(name, ID string) error {
	j.ID = ID
	return nil
}
