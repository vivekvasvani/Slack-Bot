package client

import (
	"log"
	"strings"
	"github.com/NodePrime/jsonpath"
)

func GetResponse(key string, responseBody []byte) string {
	//log.Println(key)
	var value string
	paths, err := jsonpath.ParsePaths(key)
	if err != nil {
		log.Println("1",err)
		return ""
	}
	eval, err1 := jsonpath.EvalPathsInBytes(responseBody, paths)
	if err1 != nil {
		log.Println("1",err1)
		return ""
	}
	for {
		if result, ok := eval.Next(); ok {
			value = strings.TrimSpace(result.Pretty(false))
		} else {
			break
		}
	}
	if eval.Error != nil {
		log.Println(eval.Error)
		return ""
	}
	//log.Println(value)
	return value
}

func GetResponseKeyValueAsSlice(key string, responseBody []byte) []string {
	//log.Println(key)
	vals := make([]string, 0)
	paths, err := jsonpath.ParsePaths(key)
	if err != nil {
		log.Println("1",err)
	}
	eval, err := jsonpath.EvalPathsInBytes(responseBody, paths)

	for {
		if result, ok := eval.Next(); ok {
			value := strings.TrimSpace(result.Pretty(false))
			vals=append(vals,value)
		} else {
			break
		}
	}
	if eval.Error != nil {
		log.Println(eval.Error)
	}
	return vals
}