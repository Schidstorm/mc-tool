package job

import (
	"mc-tool/zip"
	"os"
	"os/exec"
	"path"
	"runtime"
)

const (
	JdCliDistUrl = "https://github.com/kwart/jd-cli/releases/download/jd-cli-1.2.0/jd-cli-1.2.0-dist.zip"
)

func PrepareJdCli() error {
	workdir, err := os.Getwd()
	if err != nil {
		return err
	}

	jdZliZipFilePath := path.Join(workdir, "jd-cli", path.Base(JdCliDistUrl))
	err = os.MkdirAll("jd-cli", 0777)
	if err != nil {
		return err
	}
	response, size, err := HttpDownloadWithCache(JdCliDistUrl, jdZliZipFilePath)
	if err != nil {
		return err
	}
	defer response.Close()
	jdCliBuffer := new(ReadAtBuffer)
	err = CopyAllWithProgress(path.Base(JdCliDistUrl), size, jdCliBuffer, response)
	if err != nil {
		return err
	}
	jdCliZipPath := "jd-cli"
	_, err = zip.Unzip(jdCliBuffer, size, jdCliZipPath)
	if err != nil {
		return err
	}

	return nil
}

func DecompileWithJdCli(jarFilePath, decompiledSourceDir string) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	_ = os.RemoveAll(decompiledSourceDir)
	jdCliBinPath := path.Join(workingDir, "jd-cli", "jd-cli")
	if runtime.GOOS == "windows" {
		jdCliBinPath += ".bat"
	}
	err = os.MkdirAll(decompiledSourceDir, 0777)
	decompieCmd := exec.Command(jdCliBinPath, "--outputDirStructured", decompiledSourceDir, jarFilePath)
	decompieCmd.Dir = path.Dir(jdCliBinPath)
	decompieCmd.Stdout = os.Stdout
	decompieCmd.Stderr = os.Stderr
	err = ExecRunCmd(decompieCmd)
	if err != nil {
		return err
	}

	return nil
}
