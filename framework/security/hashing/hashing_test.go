package hashing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testPassword string = "MyPassword123!!##"
	testHash     string = "$2a$12$rW6Ska0DaVjTX/8sQGCp/.y7kl2RvF.9936Hmm27HyI0cJ78q1UOG"
)

func TestCheckPassword(t *testing.T) {
	isPasswordCorrect := CheckPassword(testHash, testPassword)
	assert.True(t, isPasswordCorrect, "Password was not correct")
}

func TestGenerateHash(t *testing.T) {
	hash, err := GenerateHash(testPassword)
	assert.NotNil(t, hash, "Hash is nil")
	assert.Nil(t, err, "Hashing the password did not succeed")
}
