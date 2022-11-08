package teamtailorgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCompany(t *testing.T) {
	SetUpTeamTailorTest()
	company, err := teamtailorconnection.GetCompany()
	assert.Equal(t, err, nil)
	assert.Equal(t, "dwkm7nacXJY", company.Data.ID)
}
