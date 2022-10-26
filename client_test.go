package dpo_test

import (
	"testing"

	"github.com/nndi-oss/go-dpo"
	"github.com/stretchr/testify/assert"
)

func TestCreateTokenWithNilToken(t *testing.T) {
	assert := assert.New(t)
	client := dpo.NewClient("", true)

	_, err := client.CreateToken(nil)
	assert.NotNil(err)
	assert.ErrorContains(err, "token must not be nil")
}
