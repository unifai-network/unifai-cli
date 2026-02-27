package main

import (
	"os"

	"unifai/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args[1:]))
}
