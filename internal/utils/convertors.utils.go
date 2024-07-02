package utils

import "strings"

func ConvertFromCSVToMapStringStruct(csv string) map[string]struct{} {
	separatedString := strings.Split(csv, ",")
	mapToReturn := make(map[string]struct{})
	for _, value := range separatedString {
		mapToReturn[value] = struct{}{}
	}
	return mapToReturn
}

func ConvertFromStringArrayToMapStringStruct(arr []string) map[string]struct{} {
	mapToReturn := make(map[string]struct{})
	for _, value := range arr {
		mapToReturn[value] = struct{}{}
	}
	return mapToReturn
}
