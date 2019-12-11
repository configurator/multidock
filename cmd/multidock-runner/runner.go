package main

import (
	"fmt"
	"os"

	"github.com/configurator/multidock"
	"github.com/configurator/multidock/pkg/runner"
	"github.com/configurator/multidock/pkg/types"
)

func main() {
	fmt.Printf("Multidock v%s\n", multidock.Version)
	manifest, err := types.ReadManifestFromFile("mani.fest")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	runner.Run(manifest)
}
