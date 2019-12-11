package generator

import (
	"errors"

	"github.com/configurator/multidock/pkg/docker"
	"github.com/configurator/multidock/pkg/types"
)

func autofillManifest(manifest *types.Manifest) error {
	for i := range manifest.Images {
		err := autofillManifestImage(&manifest.Images[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func autofillManifestImage(image *types.Image) error {
	if image.Id == "" {
		return errors.New("Invalid image id \"\"")
	}
	inspect, err := docker.Inspect(image.Id)
	if err != nil {
		return err
	}

	if image.RootPath == "" {
		image.RootPath = "/images/" + image.Id
	}

	if image.Workdir == "" {
		image.Workdir = inspect.Config.WorkingDir
		if image.Workdir == "" {
			image.Workdir = "/"
		}
	}

	if image.Entrypoint == nil {
		entrypoint := []string(inspect.Config.Entrypoint)
		image.Entrypoint = &entrypoint
	}

	if image.Cmd == nil {
		cmd := []string(inspect.Config.Cmd)
		image.Cmd = &cmd
	}

	if image.Env == nil {
		image.Env = &inspect.Config.Env
	}

	return nil
}
