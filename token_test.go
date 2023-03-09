package dpo_test

import (
	"testing"
	"time"

	"github.com/golang-malawi/go-dpo"
	"github.com/stretchr/testify/assert"
)

func TestAddService(t *testing.T) {
	assert := assert.New(t)

	token := &dpo.CreateTokenRequest{}
	token.AddService("X", "XYZ", time.Now())

	assert.NotNil(token.Services)
	assert.NotEmpty(token.Services)
}
