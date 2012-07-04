// Copyright 2012 Arne Roomann-Kurrik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fauxfile

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetPathError(path string, message string) error {
	return &os.PathError{
		Path: path,
		Err:  errors.New(message),
	}
}

type MockFilesystem struct {
	cwd  *MockFileInfo
	root *MockFileInfo
}

func NewMockFilesystem() *MockFilesystem {
	root := &MockFileInfo{
		file: &MockFile{
			name:       "/",
			path:       "/",
			filesystem: nil,
			mode:       os.ModeDir | 0755,
			modified:   time.Now(),
			data:       nil,
			parent:     nil,
			children:   map[string]*MockFileInfo{},
		},
	}
	mf := &MockFilesystem{
		cwd:  root,
		root: root,
	}
	root.file.filesystem = mf
	return mf
}

func (mf *MockFilesystem) getpath(path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(mf.cwd.file.path, path)
}

func (mf *MockFilesystem) resolve(path string) (*MockFileInfo, error) {
	path = mf.getpath(path)
	parts := strings.Split(path, string(filepath.Separator))
	ptr := mf.root
	for _, part := range parts {
		if part == "" {
			continue
		}
		if child := ptr.Child(part); child != nil {
			ptr = ptr.Child(part)
		} else {
			return nil, GetPathError(path, "Path does not exist")
		}
	}
	return ptr, nil
}

func (mf *MockFilesystem) exists(path string) bool {
	_, err := mf.resolve(path)
	return err == nil
}

func (mf *MockFilesystem) Chdir(dir string) error {
	fi, err := mf.resolve(dir)
	if err == nil {
		if fi.IsDir() {
			mf.cwd = fi
			return nil
		} else {
			return GetPathError(dir, "Path is not a directory")
		}
	}
	return err
}

func (mf *MockFilesystem) Mkdir(name string, perm os.FileMode) error {
	path := mf.getpath(name)
	parentpath, dirname := filepath.Split(path)
	fi, err := mf.resolve(parentpath)
	if err != nil {
		return err
	}
	if child := fi.Child(dirname); child != nil {
		return GetPathError(path, "Path already exists")
	}
	fi.file.children[dirname] = &MockFileInfo{
		file: &MockFile{
			name:       dirname,
			path:       path,
			filesystem: mf,
			mode:       perm | os.ModeDir,
			modified:   time.Now(),
			data:       nil,
			parent:     fi,
			children:   map[string]*MockFileInfo{},
		},
	}
	fi.file.modified = time.Now()
	return nil
}

func (mf *MockFilesystem) MkdirAll(path string, perm os.FileMode) error {
	path = mf.getpath(path)
	parts := strings.Split(path, string(filepath.Separator))
	base := "/"
	for _, part := range parts {
		if part == "" {
			continue
		}
		base = filepath.Join(base, part)
		if err := mf.Mkdir(base, perm); err != nil {
			return err
		}
	}
	return nil
}

func (mf *MockFilesystem) Remove(name string) error {
	fi, err := mf.resolve(name)
	if err != nil {
		return err
	}
	if len(fi.Children()) > 0 {
		return GetPathError(name, "Directory contains children")
	}
	delete(fi.Parent().Children(), fi.file.name)
	fi.Parent().file.modified = time.Now()
	return nil
}

func (mf *MockFilesystem) RemoveAll(path string) error {
	fi, err := mf.resolve(path)
	if err != nil {
		return err
	}
	delete(fi.Parent().Children(), fi.file.name)
	fi.Parent().file.modified = time.Now()
	return nil
}

func (mf *MockFilesystem) Rename(oldname string, newname string) error {
	return errors.New("Not implemented")
}

func (mf *MockFilesystem) Create(name string) (file File, err error) {
	path := mf.getpath(name)
	dir, filename := filepath.Split(path)
	fi, err := mf.resolve(dir)
	if err != nil {
		return nil, err
	}
	fi.file.children[filename] = &MockFileInfo{
		file: &MockFile{
			name:       filename,
			path:       path,
			filesystem: mf,
			mode:       0666,
			modified:   time.Now(),
			data:       nil,
			parent:     fi,
			children:   nil,
		},
	}
	fi.file.modified = time.Now()
	return fi.Child(filename).file, nil
}

func (mf *MockFilesystem) Open(name string) (file File, err error) {
	fi, err := mf.resolve(name)
	if err != nil {
		return nil, err
	}
	return fi.file, nil
}

func (mf *MockFilesystem) OpenFile(name string, flag int, perm os.FileMode) (file File, err error) {
	return nil, errors.New("Not implemented")
}

type MockFile struct {
	name       string
	path       string
	filesystem *MockFilesystem
	mode       os.FileMode
	modified   time.Time
	data       *[]byte
	parent     *MockFileInfo
	children   map[string]*MockFileInfo
}


func (mf *MockFile) Chdir() error {
	return mf.filesystem.Chdir(filepath.Dir(mf.path))
}

func (mf *MockFile) Chmod(mode os.FileMode) error {
	mf.mode = mode
	return nil
}

func (mf *MockFile) Close() error {
	return errors.New("Not implemented")
}

func (mf *MockFile) Name() string {
	return mf.name
}

func (mf *MockFile) Read(b []byte) (n int, err error) {
	return 0, errors.New("Not implemented")
}

func (mf *MockFile) ReadAt(b []byte, off int64) (n int, err error) {
	return 0, errors.New("Not implemented")
}

func (mf *MockFile) Readdir(n int) (fi []os.FileInfo, err error) {
	// TODO: Enable returning additional elements in subsequent calls.
	fi = make([]os.FileInfo, 0)
	limit := len(mf.children)
	if n > 0 {
		limit = n
	}
	if len(mf.children) < limit {
		err = io.EOF
	}
	if len(mf.children) == 0 {
		return
	}
	i := 0
	for _, child := range mf.children {
		if i == limit {
			break
		}
		fi = append(fi, child)
		i++
	}
	return
}

func (mf *MockFile) Readdirnames(n int) (names []string, err error) {
	fi, err := mf.Readdir(n)
	names = make([]string, len(fi))
	for i, f := range fi {
		names[i] = f.Name()
	}
	return names, err
}

func (mf *MockFile) Seek(offset int64, whence int) (ret int64, err error) {
	return 0, errors.New("Not implemented")
}

func (mf *MockFile) Stat() (fi os.FileInfo, err error) {
	return nil, errors.New("Not implemented")
}

func (mf *MockFile) Sync() (err error) {
	return errors.New("Not implemented")
}

func (mf *MockFile) Truncate(size int64) error {
	return errors.New("Not implemented")
}

func (mf *MockFile) Write(b []byte) (n int, err error) {
	return 0, errors.New("Not implemented")
}

func (mf *MockFile) WriteAt(b []byte, off int64) (n int, err error) {
	return 0, errors.New("Not implemented")
}

func (mf *MockFile) WriteString(s string) (ret int, err error) {
	return 0, errors.New("Not implemented")
}

type MockFileInfo struct {
	file *MockFile
}

func (mfi *MockFileInfo) Parent() *MockFileInfo {
	return mfi.file.parent
}

func (mfi *MockFileInfo) Children() map[string]*MockFileInfo {
	return mfi.file.children
}

func (mfi *MockFileInfo) Child(name string) *MockFileInfo {
	return mfi.file.children[name]
}

func (mfi *MockFileInfo) Name() string {
	return mfi.file.name
}

func (mfi *MockFileInfo) Size() int64 {
	return int64(len(*mfi.file.data))
}

func (mfi *MockFileInfo) Mode() os.FileMode {
	return mfi.file.mode
}

func (mfi *MockFileInfo) ModTime() time.Time {
	return mfi.file.modified
}

func (mfi *MockFileInfo) IsDir() bool {
	return mfi.file.mode.IsDir()
}

func (mfi *MockFileInfo) Sys() interface{} {
	return nil
}
