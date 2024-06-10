package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractBearerToken(t *testing.T) {
	_, err := extractBearerToken("")
	require.Error(t, err, "no header should produce error")

	_, err = extractBearerToken("Bearer")
	require.Error(t, err, "type with no value should produce error")

	token, err := extractBearerToken("Authorization abc123")
	require.NoError(t, err)

	assert.Equal(t, token, "abc123")
}
