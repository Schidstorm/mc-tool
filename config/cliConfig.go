package config

type CliRootConfig struct {
}

type CliMinecraftListConfig struct {
	OldAlpha bool `cli:"old_alpha" default:"false" usage:"list old_alpha version"`
	OldBeta  bool `cli:"old_beta" default:"false" usage:"list old_beta version"`
	Snapshot bool `cli:"snapshot" default:"false" usage:"list snapshot version"`
	Release  bool `cli:"release" default:"true" usage:"list release version"`
}

type CliMinecraftDownloadConfig struct {
	OutputPath   string `cli:"folder" default:"files" usage:"set output folder"`
	Version      string `cli:"version" default:"latest-release" usage:"set version"`
	Server       bool   `cli:"server" default:"true" usage:"download server files"`
	Client       bool   `cli:"client" default:"false" usage:"download client files"`
	Jar          bool   `cli:"jar" default:"true" usage:"download jar file"`
	Mappings     bool   `cli:"mappings" default:"true" usage:"download mappings file"`
	ShowProgress bool   `cli:"progress" default:"true" usage:"show or hide progress bar"`
}

type CliMinecraftDecompileConfig struct {
	Version            string `cli:"version" default:"latest-release" usage:"set version"`
	Server             bool   `cli:"server" default:"true" usage:"download server files"`
	Client             bool   `cli:"client" default:"false" usage:"download client files"`
	MCRemapperCloneDir string `cli:"mcRemapperPath" default:"MC-Remapper" usage:"directory where to clone MC-Remapper"`
	OutputDir          string `cli:"output" default:"files" usage:"decompiled source directory"`
}
