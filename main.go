package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kushsharma/selfupdate/cmd"
)

var (
	// Version is injected through linker flags
	Version = ""
	// Build date is injected through linker flags
	Build = ""
	// CommitSHA hash
	CommitSHA = ""
)

const (
	UpdateBinaryName = "main.update"
)

func main() {
	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	app := cmd.ApplicationInfo{
		Version: Version,
		Build:   Build,
		Commit:  CommitSHA,
	}
	fmt.Printf("Hello! Version %s, built at %s, located at %s\n", app.Version, app.Build, execPath)

	// upload this generated package file to target server and use it to identify if you need to download the updated binary
	// after fetching from server, we found the version is different from what we have, download the binary
	// else exit the update process. Its better if we can run this in a goroutine which keeps on polling ever X mins
	cmd.GeneratePackageInfo(app, execPath)

	// download this from internet or something and pass the downloaded path
	updatedBinary := fmt.Sprintf("%s/%s", filepath.Dir(execPath), UpdateBinaryName)

	err = cmd.UpdatePackage(execPath, updatedBinary, "")
	if err != nil {
		panic(err)
	}
}
