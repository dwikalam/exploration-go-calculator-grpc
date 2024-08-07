package main

import (
	"context"
	"log"
	"os"

	"github.com/dwikalam/calcgorpc/internal/app"
)

func main() {
	ctx := context.Background()

	if err := app.Run(ctx); err != nil {
		log.Printf("running app failed: %v", err)

		os.Exit(1)
	}
}
