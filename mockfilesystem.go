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
	"path"
	"path/filepath"
	"strings"
	"bytes"
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
		name:       "/",
		filesystem: nil,
		mode:       os.ModeDir | 0755,
		modified:   time.Now(),
		data:       nil,
		parent:     nil,
		children:   map[string]*MockFileInfo{},
	}
	mf := &MockFilesystem{
		cwd:  root,
		root: root,
	}
	root.filesystem = mf
	return mf
}

func (mf *MockFilesystem) getpath(path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(mf.cwd.path(), path)
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
	fi.children[dirname] = &MockFileInfo{
		name:       dirname,
		filesystem: mf,
		mode:       perm | os.ModeDir,
		modified:   time.Now(),
		data:       new(bytes.Buffer),
		parent:     fi,
		children:   map[string]*MockFileInfo{},
	}
	fi.modified = time.Now() // Update directory's timestamp.
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
	delete(fi.Parent().Children(), fi.name)
	fi.Parent().modified = time.Now()
	return nil
}

func (mf *MockFilesystem) RemoveAll(path string) error {
	fi, err := mf.resolve(path)
	if err != nil {
		return err
	}
	delete(fi.Parent().Children(), fi.name)
	fi.Parent().modified = time.Now()
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
	fi.children[filename] = &MockFileInfo{
		name:       filename,
		filesystem: mf,
		mode:       0666,
		modified:   time.Now(),
		data:       new(bytes.Buffer),
		parent:     fi,
		children:   nil,
	}
	fi.modified = time.Now()
	f := &MockFile{
		filesystem: mf,
		closed: false,
		path: path,
	}
	return f, nil
}

func (mf *MockFilesystem) Open(name string) (file File, err error) {
	_, err = mf.resolve(name)
	if err != nil {
		return nil, err
	}
	f := &MockFile{
		filesystem: mf,
		closed: false,
		path: name,
	}
	return f, nil
}

func (mf *MockFilesystem) OpenFile(name string, flag int, perm os.FileMode) (file File, err error) {
	return nil, errors.New("Not implemented")
}

type MockFile struct {
	path       string
	closed     bool
	filesystem *MockFilesystem
}

func (mf *MockFile) stat() (mfi *MockFileInfo, err error) {
	return mf.filesystem.resolve(mf.path)
}

func (mf *MockFile) Chdir() error {
	return mf.filesystem.Chdir(filepath.Dir(mf.path))
}

func (mf *MockFile) Chmod(mode os.FileMode) error {
	var (
		mfi *MockFileInfo
		err error
	)
	if mfi, err = mf.stat(); err != nil {
		return err
	}
	mfi.mode = mode
	return nil
}

func (mf *MockFile) Close() error {
	mf.closed = true
	return nil
}

func (mf *MockFile) Name() string {
	_, filename := filepath.Split(mf.path)
	return filename
}

func (mf *MockFile) Read(b []byte) (n int, err error) {
	var mfi *MockFileInfo
	if mfi, err = mf.stat(); err != nil {
		return 0, err
	}
	return mfi.data.Read(b)
}

func (mf *MockFile) ReadAt(b []byte, off int64) (n int, err error) {
	return 0, errors.New("Not implemented")
}

func (mf *MockFile) Readdir(n int) (fi []os.FileInfo, err error) {
	// TODO: Enable returning additional elements in subsequent calls.
	var (
		mfi *MockFileInfo
	)
	fi = make([]os.FileInfo, 0)
	if mfi, err = mf.stat(); err != nil {
		return nil, err
	}
	limit := len(mfi.children)
	if n > 0 {
		limit = n
	}
	if len(mfi.children) < limit {
		err = io.EOF
	}
	if len(mfi.children) == 0 {
		return
	}
	i := 0
	for _, child := range mfi.children {
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
	return mf.stat()
}

func (mf *MockFile) Sync() (err error) {
	return nil
}

func (mf *MockFile) Truncate(size int64) error {
	var mfi *MockFileInfo
	if mfi, err = mf.stat(); err != nil {
		return 0, err
	}
	return mfi.data.Truncate(size)
}

func (mf *MockFile) Write(b []byte) (n int, err error) {
	var mfi *MockFileInfo
	if mfi, err = mf.stat(); err != nil {
		return 0, err
	}
	return mfi.data.Write(b)
}

func (mf *MockFile) WriteAt(b []byte, off int64) (n int, err error) {
	return 0, errors.New("Not implemented")
}

func (mf *MockFile) WriteString(s string) (ret int, err error) {
	var mfi *MockFileInfo
	if mfi, err = mf.stat(); err != nil {
		return 0, err
	}
	return mfi.data.Write(s)
}

type MockFileInfo struct {
	data       *bytes.Buffer
	name       string
	filesystem *MockFilesystem
	mode       os.FileMode
	modified   time.Time
	parent     *MockFileInfo
	children   map[string]*MockFileInfo
}

func (mfi *MockFileInfo) path() string {
	ptr := mfi
	filepath := mfi.name
	for ptr != mfi.filesystem.root {
		ptr = ptr.parent
		filepath = path.Join(ptr.name, filepath)
	}
	return filepath
}

func (mfi *MockFileInfo) Parent() *MockFileInfo {
	return mfi.parent
}

func (mfi *MockFileInfo) Children() map[string]*MockFileInfo {
	return mfi.children
}

func (mfi *MockFileInfo) Child(name string) *MockFileInfo {
	return mfi.children[name]
}

func (mfi *MockFileInfo) Name() string {
	return mfi.name
}

func (mfi *MockFileInfo) Size() int64 {
	return int64(mfi.data.Len())
}

func (mfi *MockFileInfo) Mode() os.FileMode {
	return mfi.mode
}

func (mfi *MockFileInfo) ModTime() time.Time {
	return mfi.modified
}

func (mfi *MockFileInfo) IsDir() bool {
	return mfi.mode.IsDir()
}

func (mfi *MockFileInfo) Sys() interface{} {
	return nil
}
