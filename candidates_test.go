package teamtailorgo

import (
	"strings"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetCandidate(t *testing.T) {
	SetUpTeamTailorTest()
	candidate, err := teamtailorconnection.GetCandidate(existingCandidateID)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, "", candidate.ID)
}

func TestGetCandidate_NotExistingId_ShouldGiveError(t *testing.T) {
	SetUpTeamTailorTest()
	_, err := teamtailorconnection.GetCandidate(nonExistingCandidateID)
	assert.NotEqual(t, err, nil)
	assert.Equal(t, "failure with error: The record identified by "+nonExistingCandidateID+" could not be found.", err.Error())
}

func TestGetCandidateByEmail(t *testing.T) {
	SetUpTeamTailorTest()
	candidate, err := teamtailorconnection.GetCandidateByEmail(existingCandidateEmail)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", candidate.ID)
}

func TestGetCandidateByEmail_NotExistingId_ShouldGiveError(t *testing.T) {
	SetUpTeamTailorTest()
	_, err := teamtailorconnection.GetCandidateByEmail(nonExistingCandidateEmail)
	assert.NotEqual(t, err, nil)
	assert.Equal(t, "no candidate found with email "+nonExistingCandidateEmail, err.Error())
}

func TestPostCandidate(t *testing.T) {
	SetUpTeamTailorTest()

	candidate, err := teamtailorconnection.PostCandidate(createNewUniqueCandidate())
	assert.Equal(t, err, nil)
	assert.NotEqual(t, "", candidate.ID)
}

func TestPostExistingCandidate_ShallFetchExistingCandidate(t *testing.T) {
	SetUpTeamTailorTest()

	candidate, err := teamtailorconnection.PostCandidate(createExistingCandidate())
	assert.Equal(t, err, nil)
	assert.NotEqual(t, "", candidate.ID)
}

func TestPostCandidate_BadData_ShallThrowError(t *testing.T) {
	SetUpTeamTailorTest()

	var cand CandidateRequest
	cand.Pitch = "This sounds great"
	cand.Email = ""

	_, err := teamtailorconnection.PostCandidate(cand)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "failure with error: email - can't be blank", err.Error())
}

func TestUpdateCandidate(t *testing.T) {
	SetUpTeamTailorTest()

	var candidate Candidate
	candidate.ID = "20596226"
	candidate.Firstname = "Firstname"
	candidate.Lastname = "Lastname"
	candidate.Email = existingCandidateEmailV2
	candidate.Phone = "0701234567"
	candidate.Picture = "https://irecommend.ams3.digitaloceanspaces.com/images/irecommend.png"
	candidate.Connected = true
	candidate.LinkedinURL = "linkedin.com"
	candidate.Pitch = "This sounds great"

	err := teamtailorconnection.UpdateCandidate(candidate)
	assert.Equal(t, nil, err)
}

func TestUpdateCandidate_InvalidId_ShallGiveError(t *testing.T) {
	SetUpTeamTailorTest()

	var candidate Candidate
	candidate.ID = invalidCandidateID
	candidate.Firstname = "Firstname"
	candidate.Lastname = "Lastname"
	candidate.Email = existingCandidateEmailV2
	candidate.Phone = "0701234567"
	candidate.Picture = "https://irecommend.ams3.digitaloceanspaces.com/images/irecommend.png"
	candidate.Connected = true
	candidate.LinkedinURL = "linkedin.com"
	candidate.Pitch = "This sounds great"

	err := teamtailorconnection.UpdateCandidate(candidate)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "failure with error: The record identified by "+invalidCandidateID+" could not be found.", err.Error())
}

func TestGetCandidates(t *testing.T) {
	SetUpTeamTailorTest()
	candidates, err := teamtailorconnection.GetCandidates()
	assert.Equal(t, err, nil)
	assert.Greater(t, len(candidates), 0)
}

func TestPostCandidateResume(t *testing.T) {
	SetUpTeamTailorTest()

	var candidateResume CandidateRequestResume
	candidateResume.Firstname = "Firstname"
	candidateResume.Lastname = "Lastname"
	candidateResume.Connected = false
	candidateResume.Sourced = true
	candidateResume.Phone = "0701234567"
	candidateResume.Pitch = "This sounds great"
	candidateResume.Tags = []string{"irecommend", "referred by employee"}
	candidateResume.Created = time.Now()
	candidateResume.UpdatedAt = time.Now()

	//Create dummy email
	trimmedID := strings.Replace(uuid.NewV4().String()[:13], "-", "", -1)
	dummyEmail := trimmedID + "@irecommend.se"
	candidateResume.Email = dummyEmail

	candidateResume.Resume = publicPDFLink

	candidate, err := teamtailorconnection.PostCandidateResume(candidateResume)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, "", candidate.ID)
}

func TestPostCandidateResume_InvalidURL_ShallGiveError(t *testing.T) {
	SetUpTeamTailorTest()

	var candidateResume CandidateRequestResume
	candidateResume.Firstname = "Firstname"
	candidateResume.Lastname = "Lastname"
	candidateResume.Connected = false
	candidateResume.Sourced = true
	candidateResume.Phone = "0701234567"
	candidateResume.Pitch = "This sounds great"
	candidateResume.Tags = []string{"irecommend", "referred by employee"}
	candidateResume.Created = time.Now()
	candidateResume.UpdatedAt = time.Now()

	//Create dummy email
	trimmedID := strings.Replace(uuid.NewV4().String()[:13], "-", "", -1)
	dummyEmail := trimmedID + "@irecommend.se"
	candidateResume.Email = dummyEmail

	candidateResume.Resume = privatePDFLink

	_, err := teamtailorconnection.PostCandidateResume(candidateResume)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "failure with error: resume - is invalid", err.Error())
}
