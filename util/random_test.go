package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomString(t *testing.T) {
	randomString := RandomString(6)
	require.NotEmpty(t, randomString)
}

func TestRandomLongURL(t *testing.T) {
	longURL := RandomLongURL()
	require.NotEmpty(t, longURL)
}
