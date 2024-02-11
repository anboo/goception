package goception

import (
	"reflect"
	"testing"
)

type Suite interface {
	Before(t *testing.T)
	After(t *testing.T)
}

func RunSuites(t *testing.T, suites ...Suite) {
	for _, suite := range suites {
		suiteType := reflect.TypeOf(suite)
		tests := make([]func(*testing.T), 0)

		for i := 0; i < suiteType.NumMethod(); i++ {
			method := suiteType.Method(i)
			if isTestMethod(method) {
				tests = append(tests, func(t *testing.T) {
					suite.Before(t)
					method.Func.Call([]reflect.Value{reflect.ValueOf(suite), reflect.ValueOf(t)})
					suite.After(t)
				})
			}
		}

		for _, test := range tests {
			test(t)
		}
	}
}

func isTestMethod(method reflect.Method) bool {
	return len(method.Name) > 4 && method.Name[:4] == "Test"
}
