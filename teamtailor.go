package teamtailorgo

import (
	"fmt"
	"net/http"
)

type Authorization interface {
	CheckAuthorization(h http.Handler) http.Handler
}

type TeamTailor struct {
	ApiHost       string `json:"api-host" bson:"api-host"`
	Authorization string `json:"authorization" bson:"authorization"`
	APIversion    string `json:"X-Api-Version" bson:"X-Api-Version"`
}

// Create TeamTailor instance
// TODO: Check token validity when creating TT instance
func NewTeamTailor(authToken string) (TeamTailor, error) {
	version := "20161108"
	api := "https://api.teamtailor.com/v1/"

	err := CheckAuthorization(authToken)
	if err != nil {
		fmt.Println(err)
		return TeamTailor{}, err
	}

	return TeamTailor{api, authToken, version}, nil
}

// CheckAuthorization checks token validity and if it has the correct permissions
// TODO: Also check if permissions are correct, right now it only Reads and not Write
func CheckAuthorization(token string) error {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://api.teamtailor.com/v1/departments", nil)
	req.Header.Set("Authorization", "Token token="+token)
	req.Header.Set("X-Api-Version", "20161108")
	req.Header.Set("Content-Type", "application/vnd.api+json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return UnauthorizedError(resp.StatusCode)
	}

	return nil
}
