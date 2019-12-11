package runnerEmbed

//go:generate go get -u github.com/go-bindata/go-bindata/...
//go:generate echo Building multidock-runner...
//go:generate go build -o multidock-runner github.com/configurator/multidock/cmd/multidock-runner
//go:generate echo Embedding file in multidock
//go:generate go-bindata -pkg runnerEmbed -o multidock-runner.go multidock-runner

func MultidockRunnerBinary() []byte {
	return MustAsset("multidock-runner")
}
