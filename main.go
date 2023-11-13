package main

import (
	"fmt"
	"os"

	"github.com/monopole/myrepos/internal/config"
	"github.com/monopole/myrepos/internal/file"
	"github.com/monopole/myrepos/internal/ssh"
	"github.com/monopole/myrepos/internal/tree"
	"github.com/monopole/myrepos/internal/visitor"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newCommand() *cobra.Command {
	var cfg *config.Config
	return &cobra.Command{
		Use:     "myrepos [{path/to/config/file}]",
		Short:   "Clone or rebase the repositories specified in the input file.",
		Long:    "",
		Example: "",
		Args: func(_ *cobra.Command, args []string) error {
			filePath, err := file.GetFilePath(args)
			if err != nil {
				return err
			}
			cfg, err = loadConfig(filePath)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ssh.ErrIfNoSshAgent(); err != nil {
				return err
			}
			if err := ssh.ErrIfNoSshKeys(); err != nil {
				return err
			}
			t, err := tree.MakeRootNode(cfg)
			if err != nil {
				return err
			}
			v := visitor.Cloner{}
			t.Accept(&v)
			return v.Err()
		},
		SilenceUsage: true,
	}
}

func main() {
	if err := newCommand().Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func loadConfig(p file.Path) (*config.Config, error) {
	body, err := os.ReadFile(string(p))
	if err != nil {
		return nil, fmt.Errorf("unable to read config file %q", p)
	}
	var c config.Config
	if err = yaml.Unmarshal(body, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
