package configs

import (
	"fmt"
	"os"
	"strings"
)

func GetPostgresDSN() (string, error) {
	params := make([]string, 5)

	host := os.Getenv("POSTGRES_HOSTNAME")
	if host == "" {
		return "", fmt.Errorf("POSTGRES_HOSTNAME is required")
	}
	params = append(params, fmt.Sprintf("host=%s", host))

	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		return "", fmt.Errorf("POSTGRES_USER is required")
	}
	params = append(params, fmt.Sprintf("user=%s", user))

	dbname := os.Getenv("POSTGRES_DB")
	if dbname == "" {
		return "", fmt.Errorf("POSTGRES_DB is required")
	}
	params = append(params, fmt.Sprintf("dbname=%s", dbname))

	password := os.Getenv("POSTGRES_PASSWORD")
	if password != "" {
		params = append(params, fmt.Sprintf("password=%s", password))
	}

	sslmode := os.Getenv("SSL_MODE")
	if sslmode != "" {
		sslmode = fmt.Sprintf("sslmode=%s", sslmode)
	} else {
		sslmode = "sslmode=disable"
	}
	params = append(params, sslmode)

	return strings.Join(params, " "), nil
}
