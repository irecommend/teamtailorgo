package teamtailorgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFirstJobPage_ShouldPass(t *testing.T) {
	SetUpTeamTailorTest()
	jobs, err := teamtailorconnection.GetFirstJobPage()
	assert.Equal(t, err, nil)
	assert.Greater(t, len(jobs), 0)
}

func TestGetAllJobs_ShouldPass(t *testing.T) {
	SetUpTeamTailorTest()
	jobs, err := teamtailorconnection.GetAllJobs()
	assert.Equal(t, nil, err)
	assert.Greater(t, len(jobs), 0)
}

func TestGetJob_ShouldPass(t *testing.T) {
	SetUpTeamTailorTest()
	job, err := teamtailorconnection.GetJob(existingJobId)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", job.Title)
}

func TestGetJob_NotExistingJob_ShouldGenerateError(t *testing.T) {
	SetUpTeamTailorTest()
	_, err := teamtailorconnection.GetJob(invalidJobId)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "failure with error: The record identified by "+invalidJobId+" could not be found.", err.Error())
}
