package teamtailorgo

import "net/http"

type Authorization interface {
	CheckAuthorization(h http.Handler) http.Handler
}

type TeamTailor struct {
	ApiHost       string `json:"api-host" bson:"api-host"`
	Authorization string `json:"authorization" bson:"authorization"`
	APIversion    string `json:"X-Api-Version" bson:"X-Api-Version"`
}

// Create TeamTailor instance
func NewTeamTailor(authToken string) TeamTailor {
	version := "20161108"
	api := "https://api.teamtailor.com/v1/"

	return TeamTailor{api, authToken, version}
}

// Check token validity
func (t *TeamTailor) CheckAuthorization(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
