package docker

import (
	"context"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

func createTarArchive(directory string) (io.ReadCloser, error) {
	return archive.Tar(directory, archive.Gzip)
}

func Build(directory string, tags []string) error {
	docker, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	docker.NegotiateAPIVersion(context.Background())

	buildContext, err := createTarArchive(directory)
	if err != nil {
		return err
	}
	defer buildContext.Close()

	response, err := docker.ImageBuild(context.Background(), buildContext, types.ImageBuildOptions{
		Tags: tags,
	})
	if err != nil {
		return err
	}

	output, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(output)
	if err != nil {
		return err
	}

	return nil
}
