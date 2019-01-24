package helpers_test

import (
	"goapi/helpers"
	"testing"
)

func TestGenerateSignedString(t *testing.T) {
	if ss := helpers.GenerateSignedString("johndoe", "signingkey"); len(ss) == 0 {
		t.Error("GenerateSignedString() returned empty string")
	}
}
