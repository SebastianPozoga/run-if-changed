package main

import (
	"testing"
)

func TestCleanPath(t *testing.T) {
	var DataSets = []struct {
		For      string
		Expected string
	}{
		{
			For:      "./node1/../../node3",
			Expected: "node3",
		},
		{
			For:      "./node1/././node2",
			Expected: "node1/node2",
		},
		{
			For:      "node1/node2",
			Expected: "node1/node2",
		},
	}
	for _, testData := range DataSets {
		result := cleanPath(testData.For)
		if testData.Expected != result {
			t.Errorf("For input data:\n %+v \n  expected result like:\n %+v\n  and take:\n %+v", testData.For, testData.Expected, result)
		}
	}
}
