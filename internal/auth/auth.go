package auth

import (
	"fmt"
	"os"
)

func GetToken() (string, error) {
	token := os.Getenv("ASANA_PAT")
	if token == "" {
		return "", fmt.Errorf("ASANA_PAT environment variable is required")
	}
	return token, nil
}
