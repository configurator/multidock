package generator

import (
	"io/ioutil"
	"os"
	"path"
	"text/template"

	"github.com/configurator/multidock/pkg/docker"
	runnerEmbed "github.com/configurator/multidock/pkg/runner-embed"
	"github.com/configurator/multidock/pkg/types"
)

const dockerTemplateString = `
{{range $index, $image := .Images}}
	FROM {{$image.Id}} AS source_{{$index}}
{{end}}

FROM scratch AS target

{{range $index, $image := .Images}}
	COPY --from=source_{{$index}} / {{$image.RootPath}}
{{end}}

COPY ./mani.fest ./multidock /
CMD ["/multidock"]
`

func generateDockerfile(manifest *types.Manifest, dockerFilename string) error {
	dockerfile, err := os.Create(dockerFilename)
	if err != nil {
		return err
	}
	defer dockerfile.Close()

	dockerTemplate, err := template.New("").Parse(dockerTemplateString)
	if err != nil {
		return err
	}

	err = dockerTemplate.Execute(dockerfile, manifest)
	if err != nil {
		return err
	}

	return nil
}

func generateMultidockRunner(target string) error {
	data := runnerEmbed.MultidockRunnerBinary()
	return ioutil.WriteFile(target, data, 0755)
}

func generateDockerdir(manifest *types.Manifest, dockerdir string) error {
	err := types.WriteManifestToFile(manifest, path.Join(dockerdir, "mani.fest"))
	if err != nil {
		return err
	}

	err = generateDockerfile(manifest, path.Join(dockerdir, "Dockerfile"))
	if err != nil {
		return err
	}

	err = generateMultidockRunner(path.Join(dockerdir, "multidock"))
	if err != nil {
		return err
	}

	return nil
}

func GenerateDockerImage(manifest *types.Manifest, dockerdir string, tags []string) error {
	err := autofillManifest(manifest)
	if err != nil {
		return err
	}

	err = generateDockerdir(manifest, dockerdir)
	if err != nil {
		return err
	}

	err = docker.Build(dockerdir, tags)
	if err != nil {
		return err
	}

	return nil
}
