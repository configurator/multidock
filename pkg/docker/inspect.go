package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var cache = map[string]types.ImageInspect{}

func inspectUncached(id string) (types.ImageInspect, error) {
	docker, err := client.NewEnvClient()
	if err != nil {
		return types.ImageInspect{}, err
	}
	docker.NegotiateAPIVersion(context.Background())

	result, _, err := docker.ImageInspectWithRaw(context.Background(), id)
	if err != nil {
		return types.ImageInspect{}, err
	}

	return result, err
}

func Inspect(id string) (types.ImageInspect, error) {
	result, ok := cache[id]
	if ok {
		return result, nil
	}

	result, err := inspectUncached(id)
	if err != nil {
		return types.ImageInspect{}, err
	}

	cache[id] = result
	return result, nil
}
