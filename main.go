package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mikesmitty/file-search/cmd"
)

var (
	// Build info - injected at build time
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Set build info in cmd package
	cmd.Version = version
	cmd.Commit = commit
	cmd.Date = date

	ctx := context.Background()
	if err := cmd.Execute(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
