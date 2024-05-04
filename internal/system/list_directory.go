package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DirectoryCache struct {
	currentWd string
	parentWd  string
}

func InitDirectoryCache(path string) *DirectoryCache {
	cache := &DirectoryCache{currentWd: path}
	cache.generateParentWd()
	return cache
}

func (dc *DirectoryCache) GetWd() string {
	return dc.currentWd
}

func (dc *DirectoryCache) generateParentWd() {
	if strings.HasSuffix(dc.currentWd, "/") && len(dc.currentWd) > 1 {
		dc.currentWd = strings.TrimRight(dc.currentWd, "/")
	}

	dc.parentWd = filepath.Dir(dc.currentWd)
}

func (dc *DirectoryCache) Entries() []os.FileInfo {
	entries, _ := os.ReadDir(dc.currentWd)
	list := []os.FileInfo{}

	for _, entry := range entries {
		info, _ := entry.Info()
		list = append(list, info)
	}

	return list
}

func (dc *DirectoryCache) GetEntryPath(entry string) string {
  return fmt.Sprintf("%s/%s", dc.currentWd, entry)
}

func (dc *DirectoryCache) PreviousWd() {
	dc.currentWd = dc.parentWd
	dc.generateParentWd()
}

func (dc *DirectoryCache) NextWd(path string) {
	dc.parentWd = dc.currentWd
	dc.currentWd = filepath.Join(dc.parentWd, path)
}
