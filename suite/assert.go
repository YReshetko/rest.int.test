package suite

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/YReshetko/rest.int.test/util"
	"regexp"
	"strconv"
)

type conditionType byte

const (
	eq conditionType = iota
	match
	lt
	gt
	lte
	gte
	and
	or
)

var conditionTypeNames = map[string]conditionType{
	"eq":    eq,
	"match": match,
	"lt":    lt,
	"gt":    gt,
	"lte":   lte,
	"gte":   gte,
	"and":   and,
	"or":    or,
}

var conditionNameTypes = map[conditionType]string{
	eq:    "eq",
	match: "match",
	lt:    "lt",
	gt:    "gt",
	lte:   "lte",
	gte:   "gte",
	and:   "and",
	or:    "or",
}

type Assertion struct {
	variable   string `json:"var"`
	conditions []Condition
}

type Condition struct {
	condType      conditionType
	value         string
	subConditions []Condition
}

type assertionError struct {
	failedCondition *Condition
	reason          string
	actualValue     string
	checkedVariable string
}

func (err assertionError) Error() string {
	out := "\n"
	if err.reason != "" {
		out = out + "Reason: " + err.reason + ";\n"
	}
	if err.failedCondition != nil {
		out = out + "On: " + err.failedCondition.String() + ";\n"
	}
	if err.checkedVariable != "" {
		out = out + "Scope variable name: " + err.checkedVariable + ";\n"
	}
	if err.actualValue != "" {
		out = out + "Actual value: " + err.actualValue + ";\n"
	}
	return out
}

func (cond Condition) String() string {
	out := "Condition: " + conditionNameTypes[cond.condType] + ", Expected value: " + cond.value
	if len(cond.subConditions) > 0 {
		out = out + ", Sub-conditions: {"
		for _, v := range cond.subConditions {
			out = out + v.String() + "; "
		}
		out = out[:len(out)-2] + "}"
	}
	out = out + "."
	return out
}

func (a Assertion) MarshalJSON() ([]byte, error) {
	return json.Marshal(a)
}

func (t *Assertion) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	// Save variable name into assertion
	condType, ok := m["var"]
	if !ok {
		return errors.New("Any assertion must have var field to determine what requires to check")
	}
	t.variable = condType.(string)
	delete(m, "var")
	if len(m) == 0 {
		return errors.New("Assertion doesn't contain any conditions")
	}
	conditions, err := parseConditions(m)
	if err != nil {
		return err
	}
	t.conditions = conditions
	return nil
}

func parseConditions(m map[string]interface{}) ([]Condition, error) {
	conditions := []Condition{}
	for k, v := range m {
		condType, ok := conditionTypeNames[k]
		if !ok {
			return nil, errors.New(fmt.Sprintf("Unknown condition \"%s\" met into assertion", k))
		}
		condition := Condition{}
		condition.condType = condType
		switch typ := v.(type) {
		case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64:
			condition.value = fmt.Sprintf("%v", typ)
		case []interface{}:
			for _, v := range typ {
				newMap, err := util.InterfaceToMap(v)
				if err != nil {
					return nil, err
				}
				subs, err := parseConditions(newMap)
				if err != nil {
					return nil, err
				}
				condition.subConditions = append(condition.subConditions, subs...)
			}
		case interface{}:
			newMap, err := util.InterfaceToMap(typ)
			if err != nil {
				return nil, err
			}
			subs, err := parseConditions(newMap)
			if err != nil {
				return nil, err
			}
			condition.subConditions = subs
		default:
			return nil, errors.New(fmt.Sprintf("Unknown type to parse: %v", typ))
		}
		conditions = append(conditions, condition)
	}
	return conditions, nil
}

func (a Assertion) Assert(scope map[string]string) (bool, error) {
	value, ok := scope[a.variable]
	if !ok {
		return false, &assertionError{
			reason:          "Variable was not found in current scope",
			checkedVariable: a.variable,
		}
	}
	return checkConditions(value, a.conditions)
}

func checkConditions(value string, conditions []Condition) (bool, error) {
	result := true
	for _, condition := range conditions {
		fn := condition.condType.getCondFunc()
		currentResult := fn(value, condition.value, condition.subConditions)
		result = result && currentResult
		if !currentResult {
			return result, &assertionError{
				reason:          "Assertion failed",
				actualValue:     value,
				failedCondition: &condition,
			}
		}
	}
	return result, nil
}

func (typ conditionType) getCondFunc() func(actualValue string, expectedValue string, subConditions []Condition) bool {
	switch typ {
	case eq:
		return func(actualValue string, expectedValue string, subConditions []Condition) bool {
			return actualValue == expectedValue
		}
	case match:
		return func(actualValue string, expectedValue string, subConditions []Condition) bool {
			re := regexp.MustCompile(expectedValue)
			return re.Match([]byte(actualValue))
		}
	case or:
		return func(actualValue string, expectedValue string, subConditions []Condition) bool {
			result := false
			for _, condition := range subConditions {
				interimResult, _ := checkConditions(actualValue, []Condition{condition})
				result = result || interimResult
			}
			return result
		}
	case and:
		return func(actualValue string, expectedValue string, subConditions []Condition) bool {
			result := true
			for _, condition := range subConditions {
				interimResult, _ := checkConditions(actualValue, []Condition{condition})
				result = result && interimResult
			}
			return result
		}
	default:
		var floatFn func(a, b float64) bool
		switch typ {
		case lt:
			floatFn = func(a, b float64) bool {
				return a < b
			}
		case gt:
			floatFn = func(a, b float64) bool {
				return a > b
			}
		case lte:
			floatFn = func(a, b float64) bool {
				return a <= b
			}
		case gte:
			floatFn = func(a, b float64) bool {
				return a >= b
			}
		default:
			panic(fmt.Sprintf("Unknown condition type: %d", typ))
		}
		return func(actualValue string, expectedValue string, subConditions []Condition) bool {
			actual, err := strconv.ParseFloat(actualValue, 32)
			if err != nil {
				panic("Compared non numeric types")
			}
			expected, err := strconv.ParseFloat(expectedValue, 32)
			if err != nil {
				panic("Compared non numeric types")
			}
			return floatFn(actual, expected)

		}
	}
}
