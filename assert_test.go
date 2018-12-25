package main

import (
	"fmt"
	"testing"
)

func TestAssertion_Assert_Eq(t *testing.T) {
	assert := Assertion{
		variable: "some-var",
		conditions: []Condition{
			Condition{
				condType: eq,
				value:    "OK",
			},
		},
	}
	m := map[string]string{
		"some-var": "OK",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "FAIL",
	}
	runTest(assert, m, t, false, true)
}

func TestAssertion_Assert_Lt_Int(t *testing.T) {
	assert := Assertion{
		variable: "some-var",
		conditions: []Condition{
			Condition{
				condType: lt,
				value:    "100",
			},
		},
	}
	m := map[string]string{
		"some-var": "99",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "100",
	}
	runTest(assert, m, t, false, true)
	m = map[string]string{
		"some-var": "100.01",
	}
	runTest(assert, m, t, false, true)
}

func TestAssertion_Assert_Lt_Float(t *testing.T) {
	assert := Assertion{
		variable: "some-var",
		conditions: []Condition{
			Condition{
				condType: lt,
				value:    "100.00001",
			},
		},
	}
	m := map[string]string{
		"some-var": "100.00",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "100.00001",
	}
	runTest(assert, m, t, false, true)
	m = map[string]string{
		"some-var": "100.00002",
	}
	runTest(assert, m, t, false, true)
}

func TestAssertion_Assert_Lte_Int(t *testing.T) {
	assert := Assertion{
		variable: "some-var",
		conditions: []Condition{
			Condition{
				condType: lte,
				value:    "100",
			},
		},
	}
	m := map[string]string{
		"some-var": "100",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "10",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "101",
	}
	runTest(assert, m, t, false, true)
}

func TestAssertion_Assert_Lte_Float(t *testing.T) {
	assert := Assertion{
		variable: "some-var",
		conditions: []Condition{
			Condition{
				condType: lte,
				value:    "100.0001",
			},
		},
	}
	m := map[string]string{
		"some-var": "100.0001",
	}
	runTest(assert, m, t, true, false)

	m = map[string]string{
		"some-var": "10.99999",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "100.0002",
	}
	runTest(assert, m, t, false, true)
}

func TestAssertion_Assert_Gt_Int(t *testing.T) {
	assert := Assertion{
		variable: "some-var",
		conditions: []Condition{
			Condition{
				condType: gt,
				value:    "100",
			},
		},
	}
	m := map[string]string{
		"some-var": "101",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "99",
	}
	runTest(assert, m, t, false, true)
}

func TestAssertion_Assert_Gt_Float(t *testing.T) {
	assert := Assertion{
		variable: "some-var",
		conditions: []Condition{
			Condition{
				condType: gt,
				value:    "100.00001",
			},
		},
	}
	m := map[string]string{
		"some-var": "100.00002",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "100.00001",
	}
	runTest(assert, m, t, false, true)
	m = map[string]string{
		"some-var": "100",
	}
	runTest(assert, m, t, false, true)
}

func TestAssertion_Assert_Gte_Int(t *testing.T) {
	assert := Assertion{
		variable: "some-var",
		conditions: []Condition{
			Condition{
				condType: gte,
				value:    "100",
			},
		},
	}
	m := map[string]string{
		"some-var": "100",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "101",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "99",
	}
	runTest(assert, m, t, false, true)
}

func TestAssertion_Assert_Gte_Float(t *testing.T) {
	assert := Assertion{
		variable: "some-var",
		conditions: []Condition{
			Condition{
				condType: gte,
				value:    "100.00001",
			},
		},
	}
	m := map[string]string{
		"some-var": "100.00001",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "100.00002",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "100.00000",
	}
	runTest(assert, m, t, false, true)
}

func TestAssertion_Assert_Lt_Gt(t *testing.T) {
	assert := Assertion{
		variable: "some-var",
		conditions: []Condition{
			Condition{
				condType: lt,
				value:    "100",
			},
			Condition{
				condType: gt,
				value:    "90",
			},
		},
	}
	m := map[string]string{
		"some-var": "91",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "99",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "100",
	}
	runTest(assert, m, t, false, true)
	m = map[string]string{
		"some-var": "90",
	}
	runTest(assert, m, t, false, true)
}

func TestAssertion_Assert_Lte_Gte(t *testing.T) {
	assert := Assertion{
		variable: "some-var",
		conditions: []Condition{
			Condition{
				condType: lte,
				value:    "100",
			},
			Condition{
				condType: gte,
				value:    "90",
			},
		},
	}
	m := map[string]string{
		"some-var": "91",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "99",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "100",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "90",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "101",
	}
	runTest(assert, m, t, false, true)
	m = map[string]string{
		"some-var": "89",
	}
	runTest(assert, m, t, false, true)
}

func TestAssertion_Assert_And(t *testing.T) {
	assert := Assertion{
		variable: "some-var",
		conditions: []Condition{
			Condition{
				condType: and,
				subConditions: []Condition{
					Condition{
						condType: lte,
						value:    "100",
					},
					Condition{
						condType: gte,
						value:    "90",
					},
				},
			},
		},
	}
	m := map[string]string{
		"some-var": "91",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "99",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "100",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "90",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "101",
	}
	runTest(assert, m, t, false, true)
	m = map[string]string{
		"some-var": "89",
	}
	runTest(assert, m, t, false, true)
}

func TestAssertion_Assert_Or(t *testing.T) {
	assert := Assertion{
		variable: "some-var",
		conditions: []Condition{
			Condition{
				condType: or,
				subConditions: []Condition{
					Condition{
						condType: lte,
						value:    "50",
					},
					Condition{
						condType: gte,
						value:    "100",
					},
				},
			},
		},
	}
	m := map[string]string{
		"some-var": "50",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "49",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "100",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "101",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "51",
	}
	runTest(assert, m, t, false, true)
	m = map[string]string{
		"some-var": "99",
	}
	runTest(assert, m, t, false, true)
}

func TestAssertion_Assert_Match(t *testing.T) {
	assert := Assertion{
		variable: "some-var",
		conditions: []Condition{
			Condition{
				condType: match,
				value:    "\\d{3}",
			},
		},
	}
	m := map[string]string{
		"some-var": "152",
	}
	runTest(assert, m, t, true, false)
	m = map[string]string{
		"some-var": "1b6",
	}
	runTest(assert, m, t, false, true)
}

func runTest(assertion Assertion, m map[string]string, t *testing.T, expected bool, expectError bool) {
	result, err := assertion.Assert(m)
	if !expectError && err != nil {
		t.Errorf(fmt.Sprintf("Should'n be any errors like %v during assert execution\n", err))
		t.FailNow()
	}
	if (expected && !result) || (!expected && result) {
		msg := "Result shoudl be OK"
		if !expected {
			msg = "Result shoudl NOT be OK"
		}
		t.Errorf(msg)
		t.FailNow()
	}
}
