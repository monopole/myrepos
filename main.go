package main

import (
	"fmt"
	"github.com/monopole/myrepos/internal/file"
	"github.com/monopole/myrepos/internal/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
)

var lastErr error

func newCommand() *cobra.Command {
	var repos []*pkg.ValidatedRepo

	return &cobra.Command{
		Use:     "myrepos",
		Short:   "Clone or rebase the repositories specified in the input file.",
		Long:    "",
		Example: "",
		Args: func(_ *cobra.Command, args []string) error {
			filePath, err := file.GetFilePath(args)
			if err != nil {
				return err
			}
			var cfg *pkg.MyReposConfig
			cfg, err = loadConfig(filePath)
			if err != nil {
				return err
			}
			repos, err = cfg.ToRepos()
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, vr := range repos {
				fmt.Print(vr.Title())
				exists, isDir := vr.FullDir().Exists()
				if exists && !isDir {
					reportErr(fmt.Errorf("%q exists but isn't a directory", vr.FullDir()))
					continue
				}
				var (
					err     error
					outcome pkg.Outcome
				)
				if exists {
					outcome, err = vr.Rebase()
				} else {
					outcome, err = vr.Clone()
				}
				if err != nil {
					reportErr(err)
					continue
				}
				var status string
				status, err = vr.LastLog()
				if err != nil {
					reportErr(err)
					continue
				}
				reportStatus(outcome, status)
			}
			return lastErr
		},
		SilenceUsage: true,
	}
}

func reportErr(err error) {
	lastErr = err
	reportStatus(pkg.Oops, err.Error())
}

func reportStatus(outcome pkg.Outcome, status string) {
	fmt.Printf("%20s %s\n", outcome, status)
}

func main() {
	if err := newCommand().Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func loadConfig(p file.Path) (*pkg.MyReposConfig, error) {
	body, err := os.ReadFile(string(p))
	if err != nil {
		return nil, fmt.Errorf("unable to read config file %q", p)
	}
	var c pkg.MyReposConfig
	if err = yaml.Unmarshal(body, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
