package job

import (
	"github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

func ExecRunCmd(cmd *exec.Cmd) error {
	logrus.Infof("Running '%s' with Arguments '%s'", cmd.Path, strings.Join(cmd.Args, " "))
	return cmd.Run()
}
