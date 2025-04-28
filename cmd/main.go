package main

import (
	"effective-mobile-test-task/internal/app"
	"fmt"
	"os"
)

// @title           Effective Mobile API
// @version         1.0
// @description     Тестовое задание на GO для Effective Mobile.

// @host      localhost:8080
// @BasePath  /
func main() {
	builder := app.NewAppBuilder().
		WithEnv().
		WithLogger().
		WithRouter().
		WithDatabase().
		WithMigrations().
		WithUserRouter().
		WithServer()

	if err := builder.Build(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build app: %v\n", err)
		os.Exit(1)
	}

	builder.Run()
}
