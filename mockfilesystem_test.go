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
	"testing"
)

func ExpectCwd(t *testing.T, expected string, mf *MockFilesystem) {
	if mf.cwd.file.path != expected {
		t.Fatalf("Expected cwd of '%v', got '%v'", expected, mf.cwd)
	}
}

func ExpectDir(t *testing.T, path string, mf *MockFilesystem) {
	fi, err := mf.resolve(path)
	if err != nil {
		t.Fatalf("Expected path of '%s' to be present", path)
	}
	if !fi.IsDir() {
		t.Fatalf("Expected '%v' to be directory", path)
	}
}

func ExpectFile(t *testing.T, path string, mf *MockFilesystem) *MockFileInfo {
	fi, err := mf.resolve(path)
	if err != nil {
		t.Fatalf("Expected file at '%s' to be present", path)
	}
	return fi
}

func TestChdir(t *testing.T) {
	mf := NewMockFilesystem()
	mf.Mkdir("/foo", 0755)
	ExpectCwd(t, "/", mf)
	mf.Chdir("foo")
	ExpectCwd(t, "/foo", mf)
}

func TestMkdir(t *testing.T) {
	mf := NewMockFilesystem()
	mf.Mkdir("/foo", 0755)
	mf.Mkdir("/bar", 0777)
	mf.Chdir("bar")
	mf.Mkdir("baz", 0777)
	ExpectDir(t, "/foo", mf)
	ExpectDir(t, "/bar", mf)
	ExpectDir(t, "/bar/baz", mf)
}

func TestMkdirAll(t *testing.T) {
	mf := NewMockFilesystem()
	err := mf.MkdirAll("/foo/bar/baz", 0755)
	if err != nil {
		t.Fatalf("Problem creating directories: %v", err)
	}
	ExpectDir(t, "/foo", mf)
	ExpectDir(t, "/foo/bar", mf)
	ExpectDir(t, "/foo/bar/baz", mf)
}

func TestCreate(t *testing.T) {
	mf := NewMockFilesystem()
	mf.Create("foo.txt")
	fi := ExpectFile(t, "/foo.txt", mf)
	if fi.Mode().Perm() != 0666 {
		t.Fatalf("New file perm %v, expected 0666", fi.Mode().Perm())
	}
}
