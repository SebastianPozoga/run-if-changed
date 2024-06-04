package main

import (
	"reflect"
	"testing"
)

func TestFindExArg(t *testing.T) {
	var DataSets = []struct {
		For             []string
		ExpectedBase    []string
		ExpectedExtArgs []string
		ExpectedStatus  bool
	}{
		{
			For:             []string{"run-if-changed", "--include", "mydir", "--exclude", "--exts", ".jpg", "--", "commandToRun", "--commandArg", "v1"},
			ExpectedBase:    []string{"run-if-changed", "--include", "mydir", "--exclude", "--exts", ".jpg"},
			ExpectedExtArgs: []string{"commandToRun", "--commandArg", "v1"},
			ExpectedStatus:  true,
		},
		{
			For:             []string{"run-if-changed", "--include", "mydir", "--exclude", "--exts", ".jpg"},
			ExpectedBase:    nil,
			ExpectedExtArgs: nil,
			ExpectedStatus:  false,
		},
	}
	for _, testData := range DataSets {
		base, extArgs, status := findExArg(testData.For)
		if !reflect.DeepEqual(base, testData.ExpectedBase) {
			t.Errorf("For input data:\n %+v \n  expected base like:\n %+v\n  and take:\n %+v", testData.For, testData.ExpectedBase, base)
		}
		if !reflect.DeepEqual(extArgs, testData.ExpectedExtArgs) {
			t.Errorf("For input data:\n %+v \n  expected array like:\n %+v\n  and take:\n %+v", testData.For, testData.ExpectedExtArgs, extArgs)
		}
		if !reflect.DeepEqual(status, testData.ExpectedStatus) {
			t.Errorf("For input data:\n %+v \n  expected status like:\n %+v\n  and take:\n %+v", testData.For, testData.ExpectedStatus, status)
		}
	}
}
