package teamtailorgo

import (
	"fmt"
	"net/http"
)

const (
	baseURL     = "https://api.teamtailor.com/v1/"
	apiVersion  = "20161108"
	contentType = "application/vnd.api+json"
)

type Authorization interface {
	CheckAuthorization(h http.Handler) http.Handler
}

type TeamTailor struct {
	APIHost    string `json:"api-host" bson:"api-host"`
	Token      string `json:"token" bson:"token"`
	APIversion string `json:"X-Api-Version" bson:"X-Api-Version"`
	HTTPClient *http.Client
}

// Create TeamTailor instance
// TODO: Check token validity when creating TT instance
func NewTeamTailor(authToken string) (TeamTailor, error) {

	err := CheckAuthorization(authToken)
	if err != nil {
		fmt.Println(err)
		return TeamTailor{}, err
	}

	return TeamTailor{baseURL, authToken, apiVersion, &http.Client{}}, nil
}

// CheckAuthorization checks token validity and if it has the correct permissions
// TODO: Also check if permissions are correct, right now it only Reads and not Write
func CheckAuthorization(token string) error {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", baseURL+"departments", nil)
	req.Header.Set("Authorization", "Token token="+token)
	req.Header.Set("X-Api-Version", apiVersion)
	req.Header.Set("Content-Type", contentType)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return UnauthorizedError(resp.StatusCode)
	}

	return nil
}
