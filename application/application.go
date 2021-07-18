package application

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"mc-tool/config"
	"mc-tool/job"
)

type Application struct {
	showProgressBar bool
}

func NewApplication() *Application {
	return &Application{}
}

func (p *Application) RunRoot(config *config.CliRootConfig) {

}

func (p *Application) RunMinecraftList(config *config.CliMinecraftListConfig) {
	versions, err := job.ListMinecraftVersion(config.OldAlpha, config.OldAlpha, config.Snapshot, config.Release)
	if err != nil {
		logrus.Error(err)
		return
	}

	for _, version := range versions.Versions {
		fmt.Printf("%s %s %s\n", version.Id, version.Type, version.ReleaseTime.Format("2006-01-02T15-04-05"))
	}
}

func (p *Application) RunMinecraftDownload(cfg *config.CliMinecraftDownloadConfig) {
	p.showProgressBar = cfg.ShowProgress
	if cfg.Client {
		if cfg.Jar {
			_, err := job.DownloadMinecraftVersionToFile(cfg.Version, cfg.OutputPath, job.VersionDownloadTypeClient, job.VersionDownloadFileTypeJar, cfg.ShowProgress)
			if err != nil {
				logrus.Error(err)
			}
		}
		if cfg.Mappings {
			_, err := job.DownloadMinecraftVersionToFile(cfg.Version, cfg.OutputPath, job.VersionDownloadTypeClient, job.VersionDownloadFileTypeMappings, cfg.ShowProgress)
			if err != nil {
				logrus.Error(err)
			}
		}
	}

	if cfg.Server {
		if cfg.Jar {
			_, err := job.DownloadMinecraftVersionToFile(cfg.Version, cfg.OutputPath, job.VersionDownloadTypeServer, job.VersionDownloadFileTypeJar, cfg.ShowProgress)
			if err != nil {
				logrus.Error(err)
			}
		}
		if cfg.Mappings {
			_, err := job.DownloadMinecraftVersionToFile(cfg.Version, cfg.OutputPath, job.VersionDownloadTypeServer, job.VersionDownloadFileTypeMappings, cfg.ShowProgress)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}

func (p *Application) RunMinecraftDecompile(cfg *config.CliMinecraftDecompileConfig) {
	if cfg.Server {
		err := job.DecompileVersion(cfg.Version, job.VersionDownloadTypeServer, cfg.MCRemapperCloneDir, true)
		if err != nil {
			logrus.Error(err)
			return
		}
	}

	if cfg.Client {
		err := job.DecompileVersion(cfg.Version, job.VersionDownloadTypeClient, cfg.MCRemapperCloneDir, true)
		if err != nil {
			logrus.Error(err)
			return
		}
	}

}
