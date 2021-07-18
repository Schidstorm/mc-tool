package job

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

var ErrEtagNotFound = errors.New("file cache not found")

func HttpGet(url string) (*http.Response, error) {
	logHttpGet(url)
	return http.Get(url)
}

func HttpDownloadWithCache(url, filePath string) (io.ReadCloser, int64, error) {
	etagString, err := getFileEtag(filePath)
	if err == ErrEtagNotFound {
		logrus.Warn(err)
		etagString = ""
	}

	httpRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	httpRequest.Header.Set("If-None-Match", etagString)
	response, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, 0, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotModified {
		fileStat, err := os.Stat(filePath)
		if err != nil {
			return nil, 0, err
		}

		f, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)

		return f, fileStat.Size(), err
	} else {
		etagString := response.Header.Get("ETag")
		err = ioutil.WriteFile(etagPathFromFilePath(filePath), []byte(etagString), os.ModePerm)
		if err != nil {
			logrus.Warn(err)
		}

		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			return nil, 0, err
		}
		err = CopyAllWithProgress(path.Base(filePath), response.ContentLength, f, response.Body)
		if err != nil {
			return nil, 0, err
		}
		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			return nil, 0, err
		}

		return f, response.ContentLength, nil
	}
}

func getFileEtag(downloadFilePath string) (string, error) {
	etagFilePath := etagPathFromFilePath(downloadFilePath)
	_, etagStatError := os.Stat(etagFilePath)
	_, fileStatError := os.Stat(downloadFilePath)
	if os.IsNotExist(etagStatError) || os.IsNotExist(fileStatError) {
		return "", ErrEtagNotFound
	}

	etagBytes, err := ioutil.ReadFile(etagFilePath)
	if err == nil || os.IsExist(fileStatError) {
		return string(etagBytes), nil
	}

	return "", err
}

func etagPathFromFilePath(filePath string) string {
	return fmt.Sprintf("%s.ecache", filePath)
}

func logHttpGet(url string) {
	logrus.Infof("GET: %s", url)
}
