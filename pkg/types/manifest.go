package types

import (
	"encoding/json"
	"io/ioutil"
)

type Manifest struct {
	Images []Image `json:"images"`
}

type Image struct {
	Id         string    `json:"id"`
	RootPath   string    `json:"rootpath"`
	Workdir    string    `json:"workdir"`
	Entrypoint *[]string `json:"entrypoint"`
	Cmd        *[]string `json:"cmd"`
	Env        *[]string `json:"env"`
}

func ReadManifest(data []byte) (*Manifest, error) {
	var result Manifest
	err := json.Unmarshal(data, &result)
	return &result, err
}

func ReadManifestFromFile(manifestFile string) (*Manifest, error) {
	json, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		return nil, err
	}

	return ReadManifest(json)
}

func WriteManifest(manifest *Manifest) (data []byte, err error) {
	return json.Marshal(manifest)
}

func WriteManifestToFile(manifest *Manifest, manifestFile string) error {
	json, err := WriteManifest(manifest)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(manifestFile, json, 0644)
}

func (image *Image) FullCommand() []string {
	entrypoint := image.Entrypoint
	if entrypoint == nil {
		entrypoint = &[]string{}
	}

	cmd := image.Cmd
	if cmd == nil {
		cmd = &[]string{}
	}

	return append(*entrypoint, (*cmd)...)
}
