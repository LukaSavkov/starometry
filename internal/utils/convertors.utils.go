package utils

import "strings"

func ConvertFromCSVToStringArray(csv string) []string {
	return strings.Split(csv, ",")
}

func ConvertFromStringArrayToPromQLQuery(queryArray []string) string {
	return strings.Join(queryArray, " or ")
}
