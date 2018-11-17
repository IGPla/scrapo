// File storage
package storage

import (
	"bytes"
	"github.com/IGPla/scrapo/logger"
	"github.com/IGPla/scrapo/tasks"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var fileStorageLogger *log.Logger

func init() {
	fileStorageLogger = logger.GetLogger("FILESTORAGE", os.Stdout)
}

type FileStorage struct {
	Prefix string
}

func (fs FileStorage) StoreData(task *tasks.Task) error {
	parsedUrl, parseErr := url.Parse(task.URL)
	if parseErr != nil {
		fileStorageLogger.Printf("Error arised parsing task.URL %v (%v)",
			task.URL,
			parseErr.Error())
		return parseErr
	}
	var basePath string
	var fileName string
	if strings.HasSuffix(parsedUrl.Path, "/") {
		basePath = filepath.Join(fs.Prefix, parsedUrl.Host, parsedUrl.Path)
		fileName = buildFilename("", parsedUrl.RawQuery, task.ContentType)

	} else if parsedUrl.Path == "" {
		basePath = filepath.Join(fs.Prefix, parsedUrl.Host)
		fileName = buildFilename("", parsedUrl.RawQuery, task.ContentType)
	} else {
		var parts []string = strings.Split(parsedUrl.Path, "/")
		var path string = strings.Join(parts[0:len(parts)-2], "/")
		basePath = filepath.Join(fs.Prefix, parsedUrl.Host, path)
		fileName = buildFilename(parts[len(parts)-1], parsedUrl.RawQuery, task.ContentType)
	}
	mkdirErr := os.MkdirAll(basePath, 0700)
	if mkdirErr != nil {
		fileStorageLogger.Printf("Error arised creating folder hierarchy for %v (%v)",
			task.URL,
			mkdirErr.Error())
		return mkdirErr
	}

	var filePath = filepath.Join(basePath, fileName)
	outFile, createErr := os.Create(filePath)

	if createErr != nil {
		fileStorageLogger.Printf("Error arised creating file for %v (%v)",
			task.URL,
			createErr.Error())
		return createErr
	}
	defer outFile.Close()
	_, copyErr := io.Copy(outFile, bytes.NewBuffer(task.Content))
	if copyErr != nil {
		fileStorageLogger.Printf("Error arised copying file for %v (%v)",
			task.URL,
			copyErr.Error())
		return copyErr
	}
	return nil
}
