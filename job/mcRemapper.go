package job

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

const (
	McRemapperRepositoryUrl      = "https://github.com/HeartPattern/MC-Remapper.git"
	McRemapperRelativeBinaryPath = "build/install/MC-Remapper/bin/MC-Remapper"
)

func PrepareMcRemapper(mcRemapperDirectory string) error {
	if err := GitCleanCloneOrPull(McRemapperRepositoryUrl, mcRemapperDirectory); err != nil {
		return err
	}

	if err := buildMcRemapper(mcRemapperDirectory); err != nil {
		return err
	}

	return nil
}

func buildMcRemapper(mcRemapperDir string) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	gradlewPath := path.Join(workingDir, mcRemapperDir, "gradlew")
	if runtime.GOOS == "windows" {
		gradlewPath = gradlewPath + ".bat"
	}
	gradlewCmd := exec.Command(gradlewPath, "installDist")
	gradlewCmd.Dir = path.Dir(gradlewPath)
	gradlewCmd.Stdout = os.Stdout
	gradlewCmd.Stderr = os.Stderr
	err = ExecRunCmd(gradlewCmd)
	if err != nil {
		return err
	}

	return nil
}

func DeobfuscateWithMcRemapper(mcRemapperDir, jarPath, mappingsPath string) (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	mcRemapperBinPath := path.Join(workingDir, mcRemapperDir, McRemapperRelativeBinaryPath)
	if runtime.GOOS == "windows" {
		mcRemapperBinPath = mcRemapperBinPath + ".bat"
	}
	deobfuscationPath := fmt.Sprintf("%s-deobfuscation.jar", strings.TrimSuffix(jarPath, ".jar"))
	mcRemapperCmd := exec.Command(mcRemapperBinPath, jarPath, mappingsPath, "--output", deobfuscationPath)
	mcRemapperCmd.Dir = path.Dir(mcRemapperBinPath)
	mcRemapperCmd.Stdout = os.Stdout
	mcRemapperCmd.Stderr = os.Stderr
	err = ExecRunCmd(mcRemapperCmd)
	if err != nil {
		return "", err
	}

	return deobfuscationPath, nil
}
