package teamtailorgo

import (
	"encoding/json"
	"io"
	"net/http"
)

type Company struct {
	Data struct {
		ID    string `json:"id"`
		Type  string `json:"type"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
		Attributes struct {
			Locale  string `json:"locale"`
			Website string `json:"website"`
			Name    string `json:"name"`
		} `json:"attributes"`
		Relationships struct {
			Manager struct {
				Links struct {
					Self    string `json:"self"`
					Related string `json:"related"`
				} `json:"links"`
			} `json:"manager"`
		} `json:"relationships"`
	} `json:"data"`
}

func (t TeamTailor) GetCompany() (Company, error) {

	var company Company

	req, _ := http.NewRequest("GET", baseURL+"company", nil)
	req.Header.Set("Authorization", "Token token="+t.Token)
	req.Header.Set("X-Api-Version", apiVersion)
	req.Header.Set("Content-Type", contentType)

	resp, err := t.HTTPClient.Do(req)
	if err != nil {
		return company, err
	}
	defer resp.Body.Close()

	err = verifyResponse(resp)
	if err != nil {
		return company, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return company, err
	}

	err = json.Unmarshal(body, &company)
	if err != nil {
		return company, err
	}

	return company, nil
}
