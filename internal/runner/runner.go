package runner

import (
	"fmt"
	"github.com/monopole/myrepos/internal/file"
	"os/exec"
	"strings"
	"time"
)

// Runner runs some program (it wraps exec).
type Runner struct {
	// programPath is the program to run.
	programPath string
	// workDir is the working directory of the subprocess.
	workDir string
	// duration is how long the subprocess gets to run.
	duration time.Duration
	// output holds both stdout and stderr from the subprocess.
	output []byte
	// errAbbrevs maps error substrings to short, one-liner errors.
	errAbbrevs map[string]string
}

// NewRunner returns a Runner if it can find the program.
func NewRunner(
	programShortName string,
	timeout time.Duration,
	errAbbrevs map[string]string) (*Runner, error) {
	p, err := exec.LookPath(programShortName)
	if err != nil {
		return nil, fmt.Errorf(
			"no executable named %q on path: %w", programShortName, err)
	}
	return &Runner{
		programPath: p,
		duration:    timeout,
		errAbbrevs:  errAbbrevs,
	}, nil
}

// SetPwd changes the workDir in which programs will run.
func (r *Runner) SetPwd(d file.Path) {
	r.workDir = string(d)
}

// GetOutput returns the combined stdout, stderr output of the most recent run.
func (r *Runner) GetOutput() string {
	return string(r.output)
}

func (r *Runner) commonError(err error) (string, bool) {
	for k, v := range r.errAbbrevs {
		if strings.Contains(err.Error(), k) || strings.Contains(r.GetOutput(), k) {
			return v, true
		}
	}
	return "", false
}

// Run a command in the dir with a timeout.
func (r *Runner) Run(args ...string) error {
	//nolint: gosec
	cmd := exec.Command(r.programPath, args...)
	if r.workDir != "" {
		cmd.Dir = r.workDir
	}
	var err error
	return TimedCall(
		cmd.String(),
		r.duration,
		func() error {
			r.output, err = cmd.CombinedOutput()
			if err != nil {
				if msg, isCommon := r.commonError(err); isCommon {
					return fmt.Errorf(
						"%s while running: %s", msg, cmd.String())
				}
				return fmt.Errorf(
					"in dir %s running %q %w\n%s",
					cmd.Dir, cmd.String(), err, r.GetOutput())
			}
			return err
		})
}

// TimedCall runs fn, failing if it doesn't complete in the given duration.
// The description is used in the timeout error message.
func TimedCall(description string, d time.Duration, fn func() error) error {
	done := make(chan error)
	timer := time.NewTimer(d)
	defer timer.Stop()
	go func() { done <- fn() }()
	select {
	case err := <-done:
		return err
	case <-timer.C:
		return NewErrTimeOut(d, description)
	}
}
