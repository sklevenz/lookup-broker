package landscape

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	LANDSCAPES string = `
	{
		"cf-eu10": {
			"cloudcontroller": "https://api.cf.eu10.hana.ondemand.com",
			"uaa": "https://uaa.cf.eu10.hana.ondemand.com",
			"labels": [
				"master",
				"aws"
				]
		},
		"cf-eu10-001": {
			"cloudcontroller": "https://api.cf.eu10-001.hana.ondemand.com",
			"uaa": "https://uaa.cf.eu10-001.hana.ondemand.com",
			"labels": [
				"scaleout",
				"aws"
				]
		},
		"cf-eu10-002": {
			"cloudcontroller": "https://api.cf.eu10-002.hana.ondemand.com",
			"uaa": "https://uaa.cf.eu10-002.hana.ondemand.com",
			"labels": [
				"scaleout",
				"aws"
				]
		}
	}
  `
)

func TestGetData(t *testing.T) {
	os.Setenv("LANDSCAPES", LANDSCAPES)

	data := Get()
	assert.Equal(t, 3, len(data))
	assert.Equal(t, "https://api.cf.eu10-002.hana.ondemand.com", data["cf-eu10-002"].CloudController)
	assert.Equal(t, "https://uaa.cf.eu10-002.hana.ondemand.com", data["cf-eu10-002"].Uaa)
	assert.Equal(t, "scaleout", data["cf-eu10-002"].Labels[0])
	assert.Equal(t, "aws", data["cf-eu10-002"].Labels[1])
}

func TestGetNoData(t *testing.T) {
	os.Unsetenv("LANDSCAPES")

	data := Get()
	assert.Equal(t, 0, len(data))
}

func TestGetWrongData(t *testing.T) {
	os.Setenv("LANDSCAPES", "this is not json")

	data := Get()
	assert.Equal(t, 0, len(data))
}
