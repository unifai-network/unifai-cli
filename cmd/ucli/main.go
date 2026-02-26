package main

import (
	"os"

	"unifai-cli/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args[1:]))
}
