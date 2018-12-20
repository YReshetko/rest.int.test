package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"github.com/YReshetko/rest.int.test/util"
	"strings"
)

type curlRunner struct {
}

func (curlRunner)Run(command string) (head, body map[string]string){
	args := parseArgs(command)
	/*for i, arg := range args{
		fmt.Printf("Arg[%d] = %s\n", i, arg)
	}*/
	cmd := exec.Command("curl", args...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error: ", err)
		panic(err)
	}
	output := string(out)
	//fmt.Println(output)

	// TODO extract body and populate it into var scope
	head, body = parseCurlOutput(output)
	return
}

func parseArgs(argsLine string) []string {
	args := make([]string, 0)
	openedBrackets := 0
	arg := ""
	quoteOpened := false
	unaryQuoteOpened := false

	for _, symbol := range argsLine {
		switch symbol {
		case '\'':
			if unaryQuoteOpened {
				openedBrackets--
				unaryQuoteOpened = false
			} else {
				openedBrackets++
				unaryQuoteOpened = true
			}
		case '"':
			if quoteOpened {
				openedBrackets--
				if openedBrackets > 0 {
					arg = arg + string(symbol)
				}
				quoteOpened = false
			} else {
				openedBrackets++
				if openedBrackets > 1 {
					arg = arg + string(symbol)
				}
				quoteOpened = true
			}

		case ' ':
			if openedBrackets != 0 {
				arg = arg + string(symbol)
			} else {
				args = append(args, arg)
				arg = ""
			}
		default:
			arg = arg + string(symbol)
		}
	}
	args = append(args, arg)
	return args
}

func parseCurlOutput(out string) (header, body  map[string]string) {
	lines := strings.Split(out, "\r\n")
	/*for i, l := range lines{
		fmt.Printf("Line[%d]:%s\n", i, l)
	}*/
	index := 0
	if out[:4] == "HTTP" {
		header = make(map[string]string)
		for i, line := range lines {
			if line == "" {
				index = i
				break
			} else {
				if strings.Contains(line, ":") {
					key := line[:strings.Index(line, ":")]
					value := line[strings.Index(line, ":")+1:]
					header[key] = strings.TrimSpace(value)
				}

			}
		}
	}
	// TODO implement GOOD body parsing with error handling
	bodyStr := strings.TrimSpace(fmt.Sprintf("%s", lines[index:]))
	bodyInitMap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(bodyStr), &bodyInitMap); err != nil {
		panic("Can not unmarshal body")
	}
	potentialBody, err := util.Parse(bodyInitMap)
	if err != nil {
		panic("Can not parse body")
	}
	body = potentialBody
	return
}