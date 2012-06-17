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
	"os"
	"strings"
	"path/filepath"
)

type MockFilesystem struct {
	dirs map[string]os.FileMode
	files map[string][]byte
	cwd string
}

func NewMockFilesystem() *MockFilesystem {
	mode := os.ModeDir | 0777
	mf := &MockFilesystem{
		cwd: "/",
		dirs: map[string]os.FileMode{
			"/": mode,
		},
		files: map[string][]byte{},
	}
	return mf
}

func (mf *MockFilesystem) getpath(path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(mf.cwd, path)
}

func (mf *MockFilesystem) exists(path string) bool {
	if _, exists := mf.dirs[path]; exists {
		return true
	} else if _, exists := mf.files[path]; exists {
		return true
	}
	return false
}

func (mf *MockFilesystem) Chdir(dir string) error {
	dir = mf.getpath(dir)
	if _, exists := mf.dirs[dir]; exists {
		mf.cwd = dir
		return nil
	}
	return &os.PathError{Path: dir}
}

func (mf *MockFilesystem) Mkdir(name string, perm os.FileMode) error {
	perm = perm | os.ModeDir
	name = mf.getpath(name)
	if mf.exists(name) {
		return &os.PathError{Path: name}
	}
	mf.dirs[name] = perm
	return nil
}

func (mf *MockFilesystem) MkdirAll(path string, perm os.FileMode) error {
	parts := strings.Split(path, string(filepath.Separator))
	base := ""
	for _, part := range(parts) {
		base = filepath.Join(base, part)
		if err := mf.Mkdir(base, perm); err != nil {
			return err
		}
	}
	return nil
}

func (mf *MockFilesystem) Remove(name string) error {
	name = mf.getpath(name)
	if mf.exists(name) {

	}
	return &os.PathError{Path: name}
}

func (mf *MockFilesystem) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (mf *MockFilesystem) Rename(oldname string, newname string) error {
	return os.Rename(oldname, newname)
}

func (mf *MockFilesystem) Create(name string) (file *os.File, err error) {
	return os.Create(name)
}

func (mf *MockFilesystem) Open(name string) (file *os.File, err error) {
	return os.Open(name)
}

func (mf *MockFilesystem) OpenFile(name string, flag int, perm os.FileMode) (file *os.File, err error) {
	return os.OpenFile(name, flag, perm)
}

type MockFile struct {
	name string
	filesystem *MockFilesystem
	mode FileMode
	modified time.Time
	data []bytes
}



type MockFileInfo struct {
	file *MockFile
}

func (mfi *MockFileInfo) Name() string {
	return mfi.file.name
}

func (mfi *MockFileInfo) Size() int64 {
	return len(mfi.file.data)
}

func (mfi *MockFileInfo) Mode() FileMode {
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
