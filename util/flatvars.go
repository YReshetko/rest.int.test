package util

import (
	"encoding/json"
	"fmt"
)

// Assuming that at the start we have complex unknown json object which we already parsed into map[string]interface{}
// Obviously we can not simply extract deep values form inner objects
// we could use reflection or some another approach, but it's pretty easy recursively extract all vars into flat map
// For example:
//
// The json
//
//{
//	"description" : "example",
//	"host" : {
//		"address":"localhost",
//		"port":"1010",
//		"timeouts" : {
//			"eachrequest":"100",
//			"total":"10000"
//		}
//	},
//	"endpoints" : ["/user", "/account", "/images"],
//	"users" : [
//		{
//			"name":"John",
//			"surname":"Gardner"
//		},
//		{
//			"name":"Simon",
//			"surname":"Hutton"
//		}
//	]
//}
//
// will be extracted into flat map[string]string
//
// vars[description] = example
// vars[host.address] = localhost
// vars[host.port] = 1010
// vars[host.timeouts.eachrequest] = 100
// vars[host.timeouts.total] = 100
// vars[endpoints.0] = /user
// vars[endpoints.1] = /account
// vars[endpoints.2] = /images
// vars[users.0.name] = John
// vars[users.0.surname] = Gardner
// vars[users.1.name] = Simon
// vars[users.1.surname] = Hutton
//
// In this case we can simply extract the information using Dots
// NOTE: Allowed only string values as a base type of values
func Parse(m map[string]interface{}) (map[string]string, error) {
	out := make(map[string]string)
	for key, value := range m {
		switch typ := value.(type) {
		case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64:
			out[key] = fmt.Sprintf("%v", typ)
		case []interface{}:
			for i, val := range typ{
				prefixKey := fmt.Sprintf("%s.%d", key, i)
				switch tp := val.(type) {
				case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64:
					index := fmt.Sprintf("%s.%d", key, i)
					out[index] = fmt.Sprintf("%v", tp)
				default:
					if err := fillMap(tp, prefixKey, out); err != nil{
						return nil, err
					}
				}
			}
		default:
			if err := fillMap(typ, key, out); err != nil{
				return nil, err
			}

		}
	}
	return out, nil
}

func fillMap(typ interface{}, key string, out map[string]string) error {
	newMap, err := InterfaceToMap(typ)
	if err != nil {
		return err
	}
	strMap, err:= Parse(newMap)
	if err != nil {
		return err
	}
	for k, v := range strMap{
		index := fmt.Sprintf("%s.%s", key, k)
		out[index] = v
	}
	return nil
}
// Converting some interface{} to map[string]interface{} via Marshal -> json -> Unmarshal process
func InterfaceToMap(value interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(value)
	if err != nil{
		return nil, err
	}
	m := make(map[string]interface{})
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, err
}