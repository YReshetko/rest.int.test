package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"github.com/YReshetko/rest.int.test/util"
)

func main() {
	s, err := load("example.json")
	//s, err := load("test_example.json")

	if err != nil {
		panic(err)
	}

	fmt.Println(s)

	variables, err := util.Parse(s.Vars)
	if err != nil {
		fmt.Println(err)
		panic("Unknown type of json")
	}
	for k, v := range variables {
		fmt.Printf("map[%s]=%s\n", k, v)
	}

}
func load(fileName string) (*Suit, error) {
	file, ok := ioutil.ReadFile(fileName)
	if ok != nil {
		err := errors.New("Can't load" + fileName)
		return nil, err
	}
	s := new(Suit)
	err := json.Unmarshal(file, s)
	return s, err
}
