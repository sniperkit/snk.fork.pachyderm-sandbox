package util

import(
	"strings"
	"fmt"

	"github.com/satori/go.uuid"
)

func GenerateUniqueName(prefix string) string {
	unique_suffix := strings.Replace(uuid.NewV4().String(), "-", "", -1)
	unique_name := UniqueNameFromToken(prefix, unique_suffix[0:12])

	return unique_name
}

func UniqueNameFromToken(prefix string, token string) string {
	return fmt.Sprintf("%v-%v", prefix, token)
}

func NameFromUniqueName(unique_name string) string {
	tokens := strings.Split(unique_name, "-")
	return strings.Join(tokens[0:len(tokens)-1], "-")
}

