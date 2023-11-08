package test

import (
	"homework/jsonutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModifyJSON(t *testing.T) {
	inputJSON := []byte(`{"a": 123, "b": [4, 5, 6], "c": {"d": 7.5, "e": 8}}`)
	expectedJSON := []byte(`{"b": [1004, 5, 1006], "c": {"d": 7.5}}`)

	outputJSON, err := jsonutil.UpdateJSON(inputJSON)
	if err != nil {
		t.Errorf("UpdateJSON returned an error: %v", err)
	}

	assert.JSONEq(t, string(expectedJSON), string(outputJSON))
}
