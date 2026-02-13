package main

import (
	"log"

	"github.com/BigCactusLabs/codex-multipass/internal/app"
)

func main() {
	if err := app.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
