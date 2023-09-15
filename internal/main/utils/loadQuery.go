package utils

import (
	"os"
)

func LoadQuery(filename string) (string, error) {
	query, err := os.ReadFile("internal/main/sql/"+filename)
	if err != nil {
		return "", err
	}

	return string(query), nil
}