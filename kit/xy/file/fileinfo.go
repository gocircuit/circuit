// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package file

import (
	"os"
	"time"
)

// FileInfo holds meta-information about a file.
type FileInfo struct {
	SaveName    string
	SaveSize    int64
	SaveMode    os.FileMode
	SaveModTime time.Time
	SaveIsDir   bool
	SaveSys     interface{}
}

// NewFileInfoOS creates a new FileInfo structure from an os.FileInfo one.
func NewFileInfoOS(fi os.FileInfo) *FileInfo {
	return &FileInfo{
		SaveName:    fi.Name(),
		SaveSize:    fi.Size(),
		SaveMode:    fi.Mode(),
		SaveModTime: fi.ModTime(),
		SaveIsDir:   fi.IsDir(),
		SaveSys:     fi.Sys(),
	}
}

// Name returns the name of the file.
func (fi *FileInfo) Name() string {
	return fi.SaveName
}

// Size returns the size of the file.
func (fi *FileInfo) Size() int64 {
	return fi.SaveSize
}

// Mode retusn the mode of the file.
func (fi *FileInfo) Mode() os.FileMode {
	return fi.SaveMode
}

// ModTime returns the time the file was last modified.
func (fi *FileInfo) ModTime() time.Time {
	return fi.SaveModTime
}

// IsDir returns true if the file is a directory.
func (fi *FileInfo) IsDir() bool {
	return fi.SaveIsDir
}

// Sys returns any auxiliary file-related data.
func (fi *FileInfo) Sys() interface{} {
	return fi.SaveSys
}
