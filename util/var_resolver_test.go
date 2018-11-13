package util

import (
	"fmt"
	"testing"
)

func TestResolve_replacement_dont_need(t *testing.T) {
	testMap := map[string]string{
		"var-1" : "value 1",
		"var-2" : "value 2",
		"var-3" : "value 3",
	}

	actualMap, err := Resolve(testMap)
	if err != nil {
		t.Error("Unexpected error during parsing map", err)
		t.FailNow()
	}

	assertEqualMaps(testMap, actualMap, t)
}

func TestResolve_simple_replacement(t *testing.T) {
	testMap := map[string]string{
		"var-1" : "value 1 + ${var-2}",
		"var-2" : "value 2 + ${var-3}",
		"var-3" : "value 3",
	}

	expectedMap := map[string]string{
		"var-1" : "value 1 + value 2 + value 3",
		"var-2" : "value 2 + value 3",
		"var-3" : "value 3",
	}

	actualMap, err := Resolve(testMap)
	if err != nil {
		t.Error("Unexpected error during parsing map", err)
		t.FailNow()
	}

	assertEqualMaps(expectedMap, actualMap, t)
}


func TestResolve_unresolved_dependency(t *testing.T) {
	testMap := map[string]string{
		"var-1" : "value 1 + ${var-2}",
		"var-2" : "value 2 + ${var-3}",
	}

	_, err := Resolve(testMap)
	if err == nil && err.Error() == "No defined variable var-3"{
		t.Error("Expected error with unresolved dependency message")
		t.FailNow()
	}
}

func TestResolve_cyclic_dependency(t *testing.T) {
	testMap := map[string]string{
		"var-1" : "value 1 + ${var-2}",
		"var-2" : "value 2 + ${var-1}",
	}

	_, err := Resolve(testMap)
	if err == nil && err.Error() == "Cycled dependency on token, cannot resolve var: value 2 + ${var-1}"{
		t.Error("Expected error with Cycled dependency message")
		t.FailNow()
	}
}

func TestResolve_complex_replacement(t *testing.T) {
	testMap := map[string]string{
		"var-1" : "value 1 and ${var-2}",
		"var-2" : "value 2 and ${var-3}",
		"var-3" : "value 3 and ${var-4} and ${var-5} and ${var-6}",
		"var-4" : "value 4 and ${var-5}",
		"var-5" : "value 5",
		"var-6" : "value 6 and ${var-5}",
	}

	expectedMap := map[string]string{
		"var-1" : "value 1 and value 2 and value 3 and value 4 and value 5 and value 5 and value 6 and value 5",
		"var-2" : "value 2 and value 3 and value 4 and value 5 and value 5 and value 6 and value 5",
		"var-3" : "value 3 and value 4 and value 5 and value 5 and value 6 and value 5",
		"var-4" : "value 4 and value 5",
		"var-5" : "value 5",
		"var-6" : "value 6 and value 5",
	}
	actualMap, err := Resolve(testMap)
	if err != nil {
		t.Error("Unexpected error during parsing map", err)
		t.FailNow()
	}
	assertEqualMaps(expectedMap, actualMap, t)
}

func assertEqualMaps(expected, actual map[string]string, t *testing.T) {
	for k, v := range actual {
		testValue, ok := expected[k]
		if !ok {
			t.Error("No expected variable into actual map: ", k)
			t.FailNow()
		}

		if v != testValue {
			t.Error(fmt.Sprintf("\nExpected value: %s,\nActual value: %s,\nFor key:%s", testValue, v, k))
			t.FailNow()
		}
	}
}
