package main

import (
	"fmt"

	"github.com/kushsharma/selfupdate/cmd"
)

var (
	// Version is injected through linker flags
	Version = "0.0"
	// Build is injected through linker flags
	Build = "0"
	// BinaryPath to executable
	BinaryPath = "./main"
)

func main() {
	app := cmd.ApplicationInfo{
		Version: Version,
		Build:   Build,
	}
	fmt.Printf("Hello! Version %s, built at %s\n", app.Version, app.Build)
	cmd.Generate(app, BinaryPath)
}
