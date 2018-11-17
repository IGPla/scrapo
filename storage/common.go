// Common storage interface
package storage

import (
	"fmt"
	"github.com/IGPla/scrapo/tasks"
	"mime"
	"sort"
	"strings"
)

type Storage interface {
	StoreData(*tasks.Task) error
}

// Build filename from parsed filename plus querystring
func buildFilename(baseName string, rawQuery string, contentType string) string {
	var fileName string
	if strings.TrimSpace(baseName) == "" {
		extensions, error := mime.ExtensionsByType(contentType)
		if error == nil && len(extensions) > 0 {
			sort.Slice(extensions, func(i, j int) bool {
				return extensions[i] < extensions[j]
			})
			baseName = fmt.Sprintf("data%v", extensions[0])
		} else {
			baseName = "data"
		}
	}
	if strings.TrimSpace(rawQuery) != "" {
		fileName = strings.Join([]string{baseName, rawQuery}, "?")
	} else {
		fileName = baseName
	}
	return fileName
}
