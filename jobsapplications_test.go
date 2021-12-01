package teamtailorgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetJobApplicationStage_ShouldPass(t *testing.T) {
	SetUpTeamTailorTest()
	jobApplicationStage, err := teamtailorconnection.GetJobApplicationStage(existingJobApplicationId)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, "", jobApplicationStage.ID)
}

func TestGetJobApplicationStage_InvalidId_ShouldGiveError(t *testing.T) {
	SetUpTeamTailorTest()
	_, err := teamtailorconnection.GetJobApplicationStage(invalidJobApplicationId)
	assert.NotEqual(t, err, nil)
	assert.Equal(t, "failure with error: "+invalidJobApplicationId+" is not a valid value for id.", err.Error())
}

func TestGetJobApplicationStage_NotExistingId_ShouldGiveError(t *testing.T) {
	SetUpTeamTailorTest()
	_, err := teamtailorconnection.GetJobApplicationStage(nonExistingJobApplicationId)
	assert.NotEqual(t, err, nil)
	assert.Equal(t, "failure with error: The record identified by "+nonExistingJobApplicationId+" could not be found.", err.Error())
}

func TestCreateJobApplication_ShouldPass(t *testing.T) {
	SetUpTeamTailorTest()

	candidate, err := teamtailorconnection.PostCandidate(createNewUniqueCandidate())
	assert.Equal(t, nil, err)

	jobApplication, err := teamtailorconnection.CreateJobApplication(existingJobId, candidate.ID)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", jobApplication.ID)
}

func TestCreateJobApplication_JobaApplicationExists_ShouldGiveError(t *testing.T) {
	SetUpTeamTailorTest()

	_, err := teamtailorconnection.CreateJobApplication(existingJobId, existingCandidateID)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "job-application for candidate and position already exist", err.Error())
}
