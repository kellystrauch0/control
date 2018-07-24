package tiller

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
	"github.com/supergiant/supergiant/pkg/tasks"
)

type Task struct {
	runner runner.Runner
	script *template.Template
	output io.Writer
}

type Config struct {
	KubernetesVersion   string
	KubernetesConfigDir string
	RBACEnabled         bool
	EtcdHost            string
	EtcdPort            string
	PrivateIpv4         string
	ProviderString      string
	MasterHost          string
	MasterPort          string
}

func New(script *template.Template,
	outStream io.Writer, cfg *ssh.Config) (*Task, error) {
	sshRunner, err := ssh.NewRunner(cfg)

	if err != nil {
		return nil, errors.Wrap(err, "error creating ssh runner")
	}

	t := &Task{
		runner: sshRunner,
		script: script,
		output: outStream,
	}

	return t, nil
}

func (j *Task) WriteManifest(config Config) error {
	err := tasks.RunTemplate(context.Background(), j.script, j.runner, j.output, config)

	if err != nil {
		return errors.Wrap(err, "error running write certificates template as a command")
	}

	return nil
}
