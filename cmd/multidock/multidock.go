package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/configurator/multidock/pkg/generator"
	"github.com/configurator/multidock/pkg/types"
)

func main() {
	fmt.Printf("Multidock generator v1.0.0\n")

	manifestFile := flag.StringP("manifest", "m", "", "Manifest file")
	tags := flag.StringArrayP("tag", "t", []string{}, "Name and optionally a tag in the 'name:tag' format")
	flag.Parse()

	buildMode(*manifestFile, *tags)
}

func buildMode(manifestFile string, tags []string) {
	if manifestFile == "" {
		log.Fatal("Manifest is required")
	}

	manifest, err := types.ReadManifestFromFile(manifestFile)
	if err != nil {
		log.Fatal(err)
	}

	dockerdir, err := ioutil.TempDir(".", "multidock-dockerdir-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dockerdir)

	err = generator.GenerateDockerImage(manifest, dockerdir, tags)
	if err != nil {
		log.Fatal(err)
	}
}
