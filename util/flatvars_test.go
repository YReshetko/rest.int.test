package util

import (
	"encoding/json"
	"strings"
	"testing"
)

var test  =  `
              {
                		"description" : "example",
                		"host" : {
                			"address":"localhost",
                			"port":"1010",
                			"timeouts" : {
                				"eachrequest":100,
                				"total":10000.356
                			}
                		},
                		"endpoints" : ["/user", "/account", "/images"],
                		"users" : [
                			{
                				"name":"John",
                				"surname":"Gardner"
                			},
                			{
                				"name":"Simon",
                				"surname":"Hutton"
                			}
                		]
              }`
var expected = map[string]string{
	"host.timeouts.eachrequest" : "100",
	"endpoints.1" : "/account",
	"endpoints.2" : "/images",
	"users.0.name" : "John",
	"users.1.name" : "Simon",
	"users.1.surname" : "Hutton",
	"host.timeouts.total" : "10000.356",
	"host.address" : "localhost",
	"host.port" : "1010",
	"endpoints.0" : "/user",
	"users.0.surname" : "Gardner",
	"description" : "example",
}
func TestParse(t *testing.T) {
	actual := retrieveActualMap(test, t)
	assertMaps(actual, expected, t)
}

var test2  =  `
              {
				"some" : {
					"deep" : [
						{
							"ignore":true
						},
						{
							"path": {
								"inner" : {
									"array":[
										"Zero",
										"One"
									]
								}
							}
						}
					]
				}
              }`
var expected2 = map[string]string{
	"some.deep.0.ignore" : "true",
	"some.deep.1.path.inner.array.0" : "Zero",
	"some.deep.1.path.inner.array.1" : "One",
}
func TestParse2(t *testing.T) {
	actual := retrieveActualMap(test2, t)
	assertMaps(actual, expected2, t)

}

func retrieveActualMap(test string, t *testing.T) map[string]string {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(strings.TrimSpace(test)), &m); err != nil {
		t.Error("initial json has to be unmarshaled")
		t.FailNow()
	}
	actual, err := Parse(m)
	if err != nil {
		t.Errorf("Should'n be any errors like %s during parsing\n", err.Error())
		t.FailNow()
	}
	return actual
}

func assertMaps(actual, expected map[string]string, t *testing.T)  {
	if len(actual) != len(expected) {
		t.Error("Expected and actual map size doesn't match")
		t.FailNow()
	}

	for k, v := range actual {
		if expected[k] != v {
			t.Errorf("Failed resulting map:\nExpected:map[%s]=%s;\nActual:map[%s]=%s\n", k, expected[k], k, v)
		}
	}
}