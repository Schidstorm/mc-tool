package job

import (
	"fmt"
	"os"
	"path"
)

func DecompileVersion(version string, downloadType VersionDownloadType, mcRemapperDir string, showProgressBar bool) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	//Download Minecraft jar and mappings
	filesFolderPath := path.Join(workingDir, "files")
	jarFilePath, err := DownloadMinecraftVersionToFile(version, filesFolderPath, downloadType, VersionDownloadFileTypeJar, showProgressBar)
	if err != nil {
		return err
	}
	mappingFilePath, err := DownloadMinecraftVersionToFile(version, filesFolderPath, downloadType, VersionDownloadFileTypeMappings, showProgressBar)
	if err != nil {
		return err
	}

	if err := PrepareMcRemapper(mcRemapperDir); err != nil {
		return err
	}

	var deobfuscatedJarPath = ""
	if deobfuscationPath, err := DeobfuscateWithMcRemapper(mcRemapperDir, jarFilePath, mappingFilePath); err != nil {
		return err
	} else {
		deobfuscatedJarPath = deobfuscationPath
	}

	//download jd-cli java decompiler from github
	if err := PrepareJdCli(); err != nil {
		return err
	}

	//decompile deobfuscated jar

	sourcesDirPath := path.Join(workingDir, "files", fmt.Sprintf("%s-%s", version, VersionDownloadTypeString(downloadType)))
	if err := DecompileWithJdCli(deobfuscatedJarPath, sourcesDirPath); err != nil {
		return err
	}

	return nil
}
