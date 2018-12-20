package main

import (
	"encoding/json"
	"fmt"
	"github.com/YReshetko/rest.int.test/util"
	"strings"
)

type requestSender byte

const (
	CURL requestSender = iota
	GO
)

var requestSenderStr = map[string]requestSender{
	"CURL": CURL,
	"GO":   GO,
}

type Suit struct {
	Description string                 `json:"description,omitempty"`
	Vars        map[string]interface{} `json:"vars,omitempty"`
	Executor    requestSender          `json:"executor"`
	Tests       []Test                 `json:"tests"`
}
type Test struct {
	Lable      string      `json:"label"`
	Command    string      `json:"command"`
	Extracts   []Extract   `json:"extracts"`
	Assertions []Assertion `json:"assertions"`
}
type Extract struct {
	Header   string `json:"header,omitempty"`
	Body     string `json:"body,omitempty"`
	Variable string `json:"var"`
}

type TestRunner interface {
	Run(command string) (head, body map[string]string)
}


func (a requestSender) MarshalJSON() ([]byte, error) {
	return json.Marshal(a)
}

func (t *requestSender) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err!=nil{
		return err
	}
	*t = requestSenderStr[str]
	return nil
}

func (s requestSender)GetRunner() TestRunner {
	switch s {
	case CURL:
		return curlRunner{}
	case GO:
		panic("GO runner has not been implemented yet")
	default:
		return curlRunner{}
	}
}

func (s Suit) Run() {
	if s.Tests == nil || len(s.Tests) == 0 {
		panic("Test suit doesn't contain any int-test")
	}

	scope, err := util.Parse(s.Vars)
	if err != nil{
		panic(err)
	}
	scope, err = util.Resolve(scope)
	if err != nil {
		panic(err)
	}
	commandRunner := s.Executor.GetRunner()
	for _, test := range s.Tests {
		command := filterStringWithTokens(test.Command, scope)
		head, body := commandRunner.Run(command)
		extracts := test.Extracts
		for _, extract := range extracts {
			if head != nil {
				extract.process(scope, head, body)
			}
		}

		assertions := test.Assertions
		for _, assertion := range assertions {
			ok, err := assertion.Assert(scope)
			if err != nil {
				panic(err)
			}
			if !ok {
				fmt.Println("Test failed: ", test.Lable)
				break
			}
		}
	}
}

func (e Extract) process(scope, head, body map[string]string) {
	if e.Header != "" {
		value := head[e.Header]
		scope[e.Variable] = value
	} else if e.Body != "" {
		value := body[e.Header]
		scope[e.Variable] = value
	} else {
		panic("Nothing to extract")
	}
}

func filterStringWithTokens(str string, v map[string]string) string {
	tokens := util.ExtractTokens(str)
	if len(tokens) == 0 {
		return str
	}
	for _, token := range tokens {
		value := util.GetValueByToken(token, v)
		str = strings.Replace(str, token, value, -1)
	}
	return str
}