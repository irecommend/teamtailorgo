package teamtailorgo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

const (
	baseURL     = "https://api.teamtailor.com/v1/"
	apiVersion  = "20211201"
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

type TeamTailorErrorResponse struct {
	Errors []struct {
		Title  string `json:"title"`
		Detail string `json:"detail"`
		Code   string `json:"code"`
		Status string `json:"status"`
	} `json:"errors"`
	Meta struct {
		Texts struct {
			Prev string `json:"prev"`
			Next string `json:"next"`
		} `json:"texts"`
	} `json:"meta"`
}

// Create TeamTailor instance
func NewTeamTailor(authToken string) (TeamTailor, error) {

	err := CheckAuthorization(authToken)
	if err != nil {
		return TeamTailor{}, err
	}

	return TeamTailor{baseURL, authToken, apiVersion, &http.Client{}}, nil
}

// CheckAuthorization checks token validity and if it has the correct permissions
// TODO: Also check if permissions are correct, right now it only checks Reads and not Write
func CheckAuthorization(token string) error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", baseURL+"departments", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Token token="+token)
	req.Header.Set("X-Api-Version", apiVersion)
	req.Header.Set("Content-Type", contentType)

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "authorization failed")
	}
	if resp.StatusCode != 200 {
		return errors.New("Request token is not valid")
	}
	defer resp.Body.Close()

	return nil
}

func verifyResponse(resp *http.Response) error {
	if !isSuccessStatusCode(resp.StatusCode) {
		if resp.StatusCode == 403 || resp.StatusCode == 401 {
			return fmt.Errorf("returning status code [%d], invalid token", resp.StatusCode)
		}

		reqBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("returning status code [%d] indicates failure but failed to read error response with error: %s", resp.StatusCode, err)
		}

		var ttErrorResponse TeamTailorErrorResponse
		err = json.Unmarshal(reqBody, &ttErrorResponse)
		if err != nil {
			return fmt.Errorf("returning status code [%d] indicates failure but failed to unmarshal error response with error: %s", resp.StatusCode, err)
		}

		return fmt.Errorf("failure with error: %s", ttErrorResponse.Errors[0].Detail)
	}
	return nil
}

func isSuccessStatusCode(code int) bool {
	return code == 200 || code == 201
}
