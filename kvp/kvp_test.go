package kvp

import (
	"reflect"
	"testing"
)

type FlippingData struct {
	Input  map[string]map[string]string
	Output map[string]map[string]bool
	Same   bool
}

var testData = []FlippingData{
	{
		Input: map[string]map[string]string{
			"A": {"FOO": "bar", "BAZ": "bum"},
			"B": {"FOO": "bar", "JAZ": "NICE"},
		},
		Output: map[string]map[string]bool{
			"FOO": {"A": true, "B": true},
			"BAZ": {"A": true, "B": false},
			"JAZ": {"A": false, "B": true},
		},
		Same: false,
	},
	{
		Input: map[string]map[string]string{
			"A": {"FOO": "bar", "BAZ": "bum"},
			"B": {"FOO": "bar", "BAZ": "bum"},
		},
		Output: map[string]map[string]bool{
			"FOO": {"A": true, "B": true},
			"BAZ": {"A": true, "B": true},
		},
		Same: true,
	},
}

func TestNameFlipping(t *testing.T) {

	for _, td := range testData {
		out, success := flipNamesWithKeys(td.Input)
		eq := reflect.DeepEqual(td.Output, out)
		if !eq || success != td.Same {
			t.Fail()
		}
	}
}
