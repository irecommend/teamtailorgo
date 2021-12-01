package teamtailorgo

import (
	"net/http"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

//Update with suitable data to run your tests
const (
	testToken                   = ""
	existingCandidateEmail      = ""
	existingCandidateEmailV2    = ""
	nonExistingCandidateEmail   = ""
	existingCandidateID         = ""
	nonExistingCandidateID      = ""
	invalidCandidateID          = ""
	publicPDFLink               = ""
	privatePDFLink              = ""
	existingJobApplicationId    = ""
	invalidJobApplicationId     = ""
	nonExistingJobApplicationId = ""
	existingJobId               = ""
	invalidJobId                = ""
)

var teamtailorconnection TeamTailor

func SetUpTeamTailorTest() {
	teamtailorconnection = TeamTailor{"https://api.teamtailor.com/v1/", testToken, "20211201", &http.Client{}}
}

func createNewUniqueCandidate() CandidateRequest {
	var cand CandidateRequest
	cand.Firstname = "Firstname"
	cand.Lastname = "Lastname"
	cand.Connected = false
	cand.Sourced = true
	cand.Pitch = "This sounds great"
	cand.Tags = []string{"irecommend", "referred by employee"}
	cand.Created = time.Now()
	cand.UpdatedAt = time.Now()

	cand.Email = generateDummyEmail()

	return cand
}

func createExistingCandidate() CandidateRequest {
	var cand CandidateRequest
	cand.Firstname = "Firstname"
	cand.Lastname = "Lastname"
	cand.Connected = false
	cand.Sourced = true
	cand.Pitch = "This sounds great"
	cand.Tags = []string{"irecommend", "referred by employee"}
	cand.Created = time.Now()
	cand.UpdatedAt = time.Now()

	cand.Email = existingCandidateEmail

	return cand
}

func generateDummyEmail() string {
	trimmedID := strings.Replace(uuid.NewV4().String()[:13], "-", "", -1)
	return trimmedID + "@irecommend.se"
}
