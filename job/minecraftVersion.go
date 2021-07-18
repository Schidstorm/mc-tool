package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type MinecraftVersionDownload struct {
	Sha1 string `json:"sha1"`
	Url  string `json:"url"`
	Size int    `json:"size"`
}

type MinecraftVersion struct {
	Downloads struct {
		Client         MinecraftVersionDownload `json:"client"`
		ClientMappings MinecraftVersionDownload `json:"client_mappings"`
		Server         MinecraftVersionDownload `json:"server"`
		ServerMappings MinecraftVersionDownload `json:"server_mappings"`
	}
}

type MinecraftVersionList struct {
	Latest struct {
		Release  string `json:"release"`
		Snapshot string `json:"snapshot"`
	} `json:"latest"`

	Versions []MinecraftVersionMeta `json:"versions"`
}

type MinecraftVersionMeta struct {
	Id          string    `json:"id"`
	Type        string    `json:"type"`
	Url         string    `json:"url"`
	Time        time.Time `json:"time"`
	ReleaseTime time.Time `json:"releaseTime"`
}

var (
	ErrorVersionNotFound = errors.New("version not found")
)

type VersionDownloadType = int

const (
	VersionDownloadTypeServer VersionDownloadType = iota
	VersionDownloadTypeClient VersionDownloadType = iota
)

const (
	MinecraftVersionManifestUrl = "https://launchermeta.mojang.com/mc/game/version_manifest.json"
)

func VersionDownloadTypeString(t VersionDownloadType) string {
	if t == VersionDownloadTypeServer {
		return "server"
	} else {
		return "client"
	}
}

type VersionDownloadFileType = int

const (
	VersionDownloadFileTypeMappings VersionDownloadFileType = iota
	VersionDownloadFileTypeJar      VersionDownloadFileType = iota
)

func ListMinecraftVersion(oldAlpha, oldBeta, snapshot, release bool) (*MinecraftVersionList, error) {
	response, err := HttpGet(MinecraftVersionManifestUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	jsonBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	versions := &MinecraftVersionList{}
	err = json.Unmarshal(jsonBody, versions)
	if err != nil {
		return nil, err
	}

	oldVersionList := versions.Versions
	versions.Versions = make([]MinecraftVersionMeta, 0)
	for _, version := range oldVersionList {
		if oldAlpha && version.Type == "old_alpha" ||
			oldBeta && version.Type == "old_beta" ||
			snapshot && version.Type == "snapshot" ||
			release && version.Type == "release" {
			versions.Versions = append(versions.Versions, version)
		}
	}

	return versions, nil
}

func DownloadMinecraftVersionToFile(version, outputFolder string, downloadType VersionDownloadType, fileType VersionDownloadFileType, showProgressBar bool) (string, error) {
	err := os.MkdirAll(outputFolder, 0777)
	if err != nil {
		return "", err
	}

	extension := "jar"
	if fileType == VersionDownloadFileTypeMappings {
		extension = "txt"
	}

	suffix := VersionDownloadTypeString(downloadType)

	versionMeta, err := GetMinecraftVersionMeta(version)
	if err != nil {
		return "", err
	}

	versionData, err := GetMinecraftVersion(versionMeta)
	if err != nil {
		return "", err
	}

	versionDownload := GetVersionDownloadLinkFromVersionAndTypes(versionData, downloadType, fileType)
	filename := fmt.Sprintf("%s-%s.%s", versionMeta.Id, suffix, extension)
	filePath := path.Join(outputFolder, filename)

	fileHash, err := Sha1OfFile(filePath)
	if err == nil && fileHash == versionDownload.Sha1 {
		return filePath, nil
	}

	reader, err := DownloadMinecraftVersionToMemory(versionDownload, downloadType, fileType)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	err = f.Truncate(0)
	if err != nil {
		return "", err
	}

	if showProgressBar {
		err = CopyAllWithProgress(filename, int64(versionDownload.Size), f, reader)
	} else {
		_, err = io.Copy(f, reader)
	}

	return filePath, nil
}

func DownloadMinecraftVersionToMemory(versionDownload MinecraftVersionDownload, downloadType VersionDownloadType, fileType VersionDownloadFileType) (io.ReadCloser, error) {
	fileResponse, err := HttpGet(versionDownload.Url)
	if err != nil {
		return nil, err
	}
	return fileResponse.Body, nil
}

func GetVersionDownloadLinkFromVersionAndTypes(versionData MinecraftVersion, downloadType VersionDownloadType, fileType VersionDownloadFileType) MinecraftVersionDownload {
	if downloadType == VersionDownloadTypeClient {
		if fileType == VersionDownloadFileTypeJar {
			return versionData.Downloads.Client
		} else {
			return versionData.Downloads.ClientMappings
		}
	} else {
		if fileType == VersionDownloadFileTypeJar {
			return versionData.Downloads.Server
		} else {
			return versionData.Downloads.ServerMappings
		}
	}
}

func GetMinecraftVersionMeta(version string) (MinecraftVersionMeta, error) {
	versionList, err := ListMinecraftVersion(true, true, true, true)
	if err != nil {
		return MinecraftVersionMeta{}, err
	}

	if version == "latest-release" {
		version = versionList.Latest.Release
	}

	if version == "latest-snapshot" {
		version = versionList.Latest.Snapshot
	}

	for _, v := range versionList.Versions {
		if v.Id == version {
			return v, nil
		}
	}

	return MinecraftVersionMeta{}, ErrorVersionNotFound
}

func GetMinecraftVersion(meta MinecraftVersionMeta) (MinecraftVersion, error) {
	metaResponse, err := HttpGet(meta.Url)
	if err != nil {
		return MinecraftVersion{}, err
	}
	defer metaResponse.Body.Close()

	metaBody, err := ioutil.ReadAll(metaResponse.Body)
	if err != nil {
		return MinecraftVersion{}, err
	}

	versionData := MinecraftVersion{}
	err = json.Unmarshal(metaBody, &versionData)
	if err != nil {
		return MinecraftVersion{}, err
	}

	return versionData, nil
}
