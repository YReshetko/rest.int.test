package util

import (
	"errors"
	"fmt"
	"strings"
)

type item struct {
	key string
	value     string
}

type items map[string]*item

// The function resolves tokens by var scope
// For example we have next vars:
// vars[key-1]=value 1 and ${key-2}
// vars[key-2]=value 2
// After the resolver all tokens like ${key-2} hve to be replaced with actual value from this map
// NOTE 1: cyclic dependency couldn't be resolved (function returns error)
// NOTE 2: if initially vars map doesn't contain some key which used in token the function returns error
func Resolve(initial map[string]string) (map[string]string, error) {
	i := items(make(map[string]*item))
	i.fromMap(initial)
	if err := i.resolveTokens(); err != nil {
		return nil, err
	}
	return i.toMap(), nil
}

func (i items)fromMap(value map[string]string) {
	for k, v := range value{
		i[k] = &item{
			k,
			v,
		}
	}
}
func (i items) toMap()	map[string]string {
	out := make(map[string]string)
	for k, v := range i{
		out[k] = v.value
	}
	return out
}

func (i items)resolveTokens() error {
	for _, v := range i {
		stack := []*item{}
		_, err := v.resolveTokens(stack, i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *item)resolveTokens(stack []*item, itms items) ([]*item, error){
	//Check if current token is already in stack
	for _, v := range stack {
		if v.value == i.value{
			return nil, errors.New(fmt.Sprintf("Cycled dependency on token, cannot resolve var: %s", i.value))
		}
	}
	// Push current item into stack
	stack = append(stack, i)
	// Get all tokens for the item
	tokens := ExtractTokens(i.value)
	// If we have tokens in value we must resolve them before returning from method
	if len(tokens) != 0 {
		for _, token := range tokens{
			// Extract variable name by token (For example token: ${some-var}, variable name: some-var)
			tokenVar := GetVarByToken(token)
			itm, ok := itms[tokenVar]
			// If items map doesn't contain such variable we cant resolve tokens
			if !ok {
				return nil, errors.New(fmt.Sprintf("No defined variable %s", tokenVar))
			}
			// Before resolve the token for current item we should resolve tokens for dependency
			var err error
			stack, err = itm.resolveTokens(stack, itms)
			if err != nil {
				return nil, err
			}
			// If all dependency are resolved then we can resolve token for current item
			i.value = strings.Replace(i.value, token, itm.value, -1)
		}
	}

	// At the end we remove current element from stack and return it for future usage
	stack = stack[:len(stack)-1]
	return stack, nil
}