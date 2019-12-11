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

func generateManifestFromArgs(images []string) *types.Manifest {
	manifest := &types.Manifest{
		Images: make([]types.Image, len(images)),
	}
	for index, id := range images {
		manifest.Images[index].Id = id
	}
	return manifest
}

func main() {
	fmt.Printf("Multidock generator v1.0.0\n")

	manifestFile := flag.StringP("manifest", "m", "", "Manifest file")
	tags := flag.StringArrayP("tag", "t", []string{}, "Name and optionally a tag in the 'name:tag' format")
	flag.Parse()

	images := flag.Args()

	manifestSpecified := *manifestFile != ""
	imagesSpecified := len(images) != 0

	if manifestSpecified && imagesSpecified {
		fmt.Println("Specify either a manifest, or a list of images, to run; not both")
		flag.Usage()
		os.Exit(1)
	}

	var manifest *types.Manifest
	var err error

	if manifestSpecified {
		manifest, err = types.ReadManifestFromFile(*manifestFile)
		if err != nil {
			log.Fatal(err)
		}
	} else if imagesSpecified {
		manifest = generateManifestFromArgs(images)
	} else {
		fmt.Println("Specify either a manifest, or a list of images, to run")
		flag.Usage()
		os.Exit(1)
	}

	dockerdir, err := ioutil.TempDir(".", "multidock-dockerdir-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dockerdir)

	err = generator.GenerateDockerImage(manifest, dockerdir, *tags)
	if err != nil {
		log.Fatal(err)
	}
}
