package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func validEnvMap() map[string]string {
	return map[string]string{
		"DEV_ENV":     "false",
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"DB_NAME":     "postgres",
		"DB_SCHEMA":   "public",
		"DB_USER":     "postgres",
		"DB_PASSWORD": "postgres",
		"DB_APP_NAME": "userplan-service",
	}
}

func setenv(m map[string]string) {
	for key, val := range m {
		os.Setenv(key, val)
	}
}

func unsetenv(m map[string]string) {
	for key := range m {
		os.Unsetenv(key)
	}
}

func TestReadEnv_Success(t *testing.T) {
	setenv(validEnvMap())

	_, err := ReadEnv()
	assert.NoError(t, err)

	unsetenv(validEnvMap())
}

func TestReadEnv_DevEnvUnset(t *testing.T) {
	setenv(validEnvMap())

	os.Unsetenv("DEV_ENV")
	_, err := ReadEnv()
	assert.Error(t, err)

	unsetenv(validEnvMap())
}

func TestReadJson(t *testing.T) {}
