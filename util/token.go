package util

import "regexp"

func ExtractTokens(str string) []string {
	re := regexp.MustCompile("\\$\\{.+?\\}")
	return re.FindAllString(str, -1)
}

func GetValueByToken(token string, m map[string]string) string {
	varName := token[2 : len(token)-1]
	value, ok := m[varName]
	if !ok {
		panic("Token value " + token + " is not declared or seted into variables")
	}
	return value
}

func GetVarByToken(token string) string {
	return token[2 : len(token)-1]
}