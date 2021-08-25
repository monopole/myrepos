package runner

import (
	"errors"
	"fmt"
	"time"
)

type errTimeOut struct {
	duration time.Duration
	cmd      string
}

func NewErrTimeOut(d time.Duration, c string) errTimeOut {
	return errTimeOut{duration: d, cmd: c}
}

func (e errTimeOut) Error() string {
	return fmt.Sprintf("hit %s timeout running '%s'", e.duration, e.cmd)
}

func IsErrTimeout(err error) bool {
	var e *errTimeOut
	return errors.As(err, &e)
}
