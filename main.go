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

const (
	version   = "v0.2.2"
	shortHelp = "Clone or rebase the repositories specified in the input file."
)

func newCommand() *cobra.Command {
	var cfg []*config.Config
	return &cobra.Command{
		Use:   "myrepos [{configFile}]",
		Short: shortHelp,
		Long:  shortHelp + " " + version,
		Example: "  myrepos " + file.DefaultConfigFileName() + `

  If the config file argument has the default value shown above,
  then the argument can be omitted.`,
		Args: func(_ *cobra.Command, args []string) error {
			paths, err := file.GetFilePath(args)
			if err != nil {
				return err
			}
			for i := range paths {
				var c *config.Config
				c, err = loadConfig(paths[i])
				if err != nil {
					return err
				}
				cfg = append(cfg, c)
			}
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ssh.ErrIfNoSshAgent(); err != nil {
				return err
			}
			if err := ssh.ErrIfNoSshKeys(); err != nil {
				return err
			}
			for i := range cfg {
				t, err := tree.MakeRootNode(cfg[i])
				if err != nil {
					return err
				}
				v := visitor.Cloner{}
				t.Accept(&v)
				if err = v.Err(); err != nil {
					return err
				}
			}
			return nil
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
