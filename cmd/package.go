package cmd

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

// ApplicationInfo contains basic information related to application binary
type ApplicationInfo struct {
	Version   string
	Build     string
	Commit    string
	BinarySHA []byte
}

// Generate creates required json file for publish
func Generate(appInfo ApplicationInfo, binPath string) {
	outputDirFlag := flag.String("o", "dist", "Output directory for writing updates")

	var defaultPlatform string
	goos := os.Getenv("GOOS")
	goarch := os.Getenv("GOARCH")
	if goos != "" && goarch != "" {
		defaultPlatform = goos + "-" + goarch
	} else {
		defaultPlatform = runtime.GOOS + "-" + runtime.GOARCH
	}
	platformFlag := flag.String("platform", defaultPlatform,
		"Target platform in the form OS-ARCH. Defaults to running os/arch or the combination of the environment variables GOOS and GOARCH if both are set.")

	platform := *platformFlag
	outputPath := *outputDirFlag

	// create output dir if does not exists
	os.MkdirAll(outputPath, 0755)
	createPackageInfo(appInfo, binPath, platform, outputPath)
}

func generateSha256(path string) []byte {
	h := sha256.New()
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	h.Write(b)
	return h.Sum(nil)
}

func createPackageInfo(appInfo ApplicationInfo, binPath, platform, outputPath string) {

	//prepare package.json
	appInfo.BinarySHA = generateSha256(binPath)
	b, err := json.MarshalIndent(appInfo, "", "    ")
	if err != nil {
		fmt.Println("error:", err)
	}
	err = ioutil.WriteFile(filepath.Join(outputPath, platform+".json"), b, 0755)
	if err != nil {
		panic(err)
	}
}
