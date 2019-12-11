package runner

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"github.com/configurator/multidock/pkg/types"
)

func Run(manifest *types.Manifest) {
	for _, image := range manifest.Images {
		err := prepareChroot(image.RootPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, image := range manifest.Images {
		go runImage(image)
	}

	select {
	// Waits forever; one of the goroutines will call os.Exit if and when needed
	}
}

func prepareChroot(directory string) error {
	// Mount directories
	for _, source := range []string{"/dev", "/proc", "/sys"} {
		target := path.Join(directory, source)
		err := os.MkdirAll(target, 0755)
		if err != nil {
			return fmt.Errorf("prepareChroot: Cannot make directory %s: %w", target, err)
		}

		err = syscall.Mount(source, target, "", syscall.MS_BIND, "")
		if err != nil {
			return fmt.Errorf("prepareChroot: Cannot mount directory %s: %w", target, err)
		}
	}

	// Mount files
	for _, source := range []string{"/etc/hostname", "/etc/mtab", "/etc/hosts", "/etc/resolv.conf"} {
		target := path.Join(directory, source)
		err := os.MkdirAll(path.Dir(target), 0755)
		if err != nil {
			return fmt.Errorf("prepareChroot: Cannot make directory for %s: %w", target, err)
		}

		os.Remove(target) // ignore errors - if file doesn't exist

		err = ioutil.WriteFile(target, []byte{}, 0644)
		if err != nil {
			return fmt.Errorf("prepareChroot: Cannot create empty file %s: %w", target, err)
		}

		err = syscall.Mount(source, target, "", syscall.MS_BIND, "")
		// err = os.Link(source, target)
		if err != nil {
			return fmt.Errorf("prepareChroot: Cannot mount file %s: %w", target, err)
		}
	}

	return nil
}

func getEnv(env *[]string, name string) string {
	if env == nil {
		return ""
	}

	prefix := name + "="
	for _, namevalue := range *env {
		if strings.HasPrefix(namevalue, prefix) {
			return namevalue[len(prefix):]
		}
	}

	return ""
}

// Runs an image in chroot, and exits the parent process if that image process
// exits, with that process's exit code.
// This way, the first exiting subcontainer causes the entire parent container
// to exit.
func runImage(image types.Image) {
	commandArgs := image.FullCommand()
	path, err := LookPath(commandArgs[0], getEnv(image.Env, "PATH"), image.RootPath)
	if err != nil {
		log.Fatalf("Process %s could not be started: %w", image.Id, err)
	}

	command := &exec.Cmd{
		Path:   path,
		Args:   commandArgs,
		Stdin:  nil,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    *image.Env,
		Dir:    image.Workdir,
		SysProcAttr: &syscall.SysProcAttr{
			Chroot: image.RootPath,
		},
	}

	err = command.Run()
	if err == nil {
		log.Printf("Process %s exited with error code %d\n", image.Id, 0)
		os.Exit(0)
	} else {
		exitErr, _ := err.(*exec.ExitError)
		if exitErr != nil {
			log.Printf("Process %s exited with error code %d\n", image.Id, exitErr.ExitCode())
			os.Exit(exitErr.ExitCode())
		} else {
			log.Printf("Process %s exited with unknown error %e\n", image.Id, err)
			os.Exit(1)
		}
	}
}
