package cluster

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"

	"github.com/appcelerator/amp/cli"
)

// Notes
//
// This is a work in progress to flesh out how the CLI should use cluster plugins.
// This first iteration isn't going to be terribly elegant; the goal is to first
// get something quickly that works, then refactor.
//
//

// ========================================================
// Public types
// ========================================================

// supported plugin providers used by the factory function `NewClusterPlugin`
const (
	Local = "local"
	AWS   = "aws"
)

// PluginConfig is used by the factory function `NewClusterPlugin` to create a new plugin instance.
type PluginConfig struct {
	// Provider is the name of the cluster provider, such as "local" or "aws"
	Provider string
	Options  map[string]string
}

// Plugin declares the methods that all plugin providers,
// such as local and aws, must implement
type Plugin interface {
	// Provider returns the name of the provider, such as "local" or "aws"
	Provider() string

	// Run executes the plugin with the specified arguments and environment variables
	Run(c cli.Interface, args []string, env map[string]string) error
}

// ========================================================
// base plugin implementation - should never be instantiated
// ========================================================

type plugin struct {
	config PluginConfig
}

func (p *plugin) Provider() string {
	return p.config.Provider
}

func (p *plugin) Run(c cli.Interface, args []string, env map[string]string) error {
	return errors.New("Run method invoked on base plugin type")
}

// ========================================================
// local plugin implementation
// ========================================================
type localPlugin struct {
	plugin
}

func (p *localPlugin) Run(c cli.Interface, args []string, env map[string]string) error {
	return queryCluster(c, args, env)
}

// ========================================================
// aws plugin implementation
// ========================================================
type awsPlugin struct {
	plugin
}

func (p *awsPlugin) Run(c cli.Interface, args []string, env map[string]string) error {
	img := "appcelerator/amp-aws"
	return RunContainer(c, img, args, env)
}

// ========================================================
// plugin factory
// ========================================================
func NewPlugin(config PluginConfig) (Plugin, error) {
	var p Plugin

	switch config.Provider {
	case Local:
		p = &localPlugin{plugin{config: config}}
	case AWS:
		p = &awsPlugin{plugin{config: config}}
	default:
		return nil, errors.New(fmt.Sprintf("Not a valid plugin provider: %s", config.Provider))
	}

	return p, nil
}

// RunContainer starts a container using the specified image for the cluster plugin.
// Cluster plugin commands are `init`, `update`, and `destroy` (provided as the single
// `args` value). Additional arguments are supplied as environment variables in `env`, not `args`.
func RunContainer(c cli.Interface, img string, args []string, env map[string]string) error {
	dockerArgs := []string{
		"run", "-t", "--rm", "--name", "amp-cluster-plugin",
		"--network", "host",
		"-v", "/var/run/docker.sock:/var/run/docker.sock",
		// TODO: remove hardcoded path
		"-v", "/Users/tony/.aws:/root/.aws",
		"-e", "GOPATH=/go",
	}

	// make environment variables available to container
	if env != nil {
		for k, v := range env {
			dockerArgs = append(dockerArgs, "-e", fmt.Sprintf("%s=%s", k, v))
		}
	}

	// this completes the docker args
	dockerArgs = append(dockerArgs, img)

	cmd := "docker"
	args = append(dockerArgs, args...)

	proc := exec.Command(cmd, args...)

	stdout, err := proc.StdoutPipe()
	if err != nil {
		return err
	}
	outscanner := bufio.NewScanner(stdout)
	go func() {
		for outscanner.Scan() {
			c.Console().Println(outscanner.Text())
		}
	}()

	stderr, err := proc.StderrPipe()
	if err != nil {
		return err
	}
	errscanner := bufio.NewScanner(stderr)
	go func() {
		for errscanner.Scan() {
			c.Console().Println(errscanner.Text())
		}
	}()

	err = proc.Start()
	if err != nil {
		return err
	}

	err = proc.Wait()
	if err != nil {
		return err
	}

	return nil
}
