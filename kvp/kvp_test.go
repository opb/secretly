package kvp

import (
	"reflect"
	"testing"
)

func TestNameFlipping(t *testing.T) {

	input := map[string]map[string]string{
		"A": {
			"FOO": "bar",
			"BAZ": "bum",
		},
		"B": {
			"FOO": "bar",
			"JAZ": "NICE",
		},
	}

	expectedOutput := map[string]map[string]bool{
		"FOO": {
			"A": true,
			"B": true,
		},
		"BAZ": {
			"A": true,
			"B": false,
		},
		"JAZ": {
			"A": false,
			"B": true,
		},
	}

	testOutput, _ := flipNamesWithKeys(input)

	eq := reflect.DeepEqual(expectedOutput, testOutput)
	if !eq {
		t.Fail()
	}
}
