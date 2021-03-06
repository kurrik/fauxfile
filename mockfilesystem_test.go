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
	"fmt"
	"os"
	"sort"
	"testing"
)

func ExpectCwd(t *testing.T, expected string, mf *MockFilesystem) {
	if mf.cwd.path() != expected {
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

func ExpectEqual(t *testing.T, expected string, actual string) {
	if expected != actual {
		t.Fatalf("Expected '%v', got '%v'", expected, actual)
	}
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

func TestMkdirAllWithChdir(t *testing.T) {
	mf := NewMockFilesystem()
	mf.MkdirAll("/home/test", 0755)
	mf.Chdir("/home/test")
	mf.Mkdir("src", 0755)
	mf.MkdirAll("/home/test/src/static", 0755)
	ExpectDir(t, "/home/test/src/static", mf)
}

func TestCreate(t *testing.T) {
	mf := NewMockFilesystem()
	mf.Create("foo.txt")
	fi := ExpectFile(t, "/foo.txt", mf)
	if fi.Mode().Perm() != 0666 {
		t.Fatalf("New file perm %v, expected 0666", fi.Mode().Perm())
	}
}

func TestCreateSubdirectoryPath(t *testing.T) {
	mf := NewMockFilesystem()
	mf.Mkdir("foo", 0755)
	mf.Chdir("foo")
	_, err := mf.Create("foo.txt")
	if err != nil {
		t.Fatalf("Create should not throw error")
	}
	fi := ExpectFile(t, "/foo/foo.txt", mf)
	ExpectEqual(t, "/foo/foo.txt", fi.path())
	ExpectEqual(t, "foo.txt", fi.name)
}

func TestLongFormCreate(t *testing.T) {
	mf := NewMockFilesystem()
	mf.Mkdir("foo", 0755)
	_, err := mf.Create("foo/foo.txt")
	if err != nil {
		t.Fatalf("Create should not throw error")
	}
	fi := ExpectFile(t, "/foo/foo.txt", mf)
	ExpectEqual(t, "/foo/foo.txt", fi.path())
	ExpectEqual(t, "foo.txt", fi.name)
}

func TestOpen(t *testing.T) {
	mf := NewMockFilesystem()
	mf.Create("foo.txt")
	f, err := mf.Open("foo.txt")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	ExpectEqual(t, "foo.txt", f.(*MockFile).Name())
	ExpectEqual(t, "foo.txt", f.(*MockFile).path)
}

func TestRemove(t *testing.T) {
	mf := NewMockFilesystem()
	mf.Mkdir("foo", 0755)
	mf.Create("foo/foo.txt")
	err := mf.Remove("foo")
	if err == nil {
		t.Fatalf("Should not be able to remove non-empty directories")
	}
	err = mf.Remove("foo/foo.txt")
	if err != nil {
		t.Fatalf("Should be able to remove file")
	}
	_, err = mf.Open("foo/foo.txt")
	if err == nil {
		t.Fatalf("Remove did not remove file")
	}
	err = mf.Remove("foo")
	if err != nil {
		t.Fatalf("Should be able to remove empty directory")
	}
	_, err = mf.Open("foo")
	if err == nil {
		t.Fatalf("Remove did not remove directory")
	}
}

func TestRemoveAll(t *testing.T) {
	mf := NewMockFilesystem()
	mf.Mkdir("foo", 0755)
	mf.Create("foo/foo.txt")
	mf.Create("foo/bar.txt")
	mf.RemoveAll("foo")
	_, err := mf.Open("/foo/foo.txt")
	if err == nil {
		t.Fatalf("RemoveAll should remove children")
	}
	_, err = mf.Open("/foo/bar.txt")
	if err == nil {
		t.Fatalf("RemoveAll should remove children")
	}
	_, err = mf.Open("/foo")
	if err == nil {
		t.Fatalf("RemoveAll should remove target")
	}
}

func TestReaddir(t *testing.T) {
	mf := NewMockFilesystem()
	mf.Mkdir("foo", 0755)
	mf.Mkdir("foo/a", 0755)
	mf.Create("foo/b.txt")
	mf.Create("foo/c.txt")
	f, _ := mf.Open("foo")
	fi, err := f.Readdir(-1)
	if err != nil {
		t.Fatalf("Readdir should not throw error.")
	}
	files := map[string]string{}
	for _, i := range fi {
		files[i.Name()] = i.(*MockFileInfo).path()
	}
	ExpectEqual(t, "/foo/a", files["a"])
	ExpectEqual(t, "/foo/b.txt", files["b.txt"])
	ExpectEqual(t, "/foo/c.txt", files["c.txt"])
	ExpectEqual(t, "3", fmt.Sprintf("%v", len(files)))
}

func TestReaddirnames(t *testing.T) {
	mf := NewMockFilesystem()
	mf.Mkdir("foo", 0755)
	mf.Mkdir("foo/a", 0755)
	mf.Create("foo/b.txt")
	mf.Create("foo/c.txt")
	f, _ := mf.Open("foo")
	names, err := f.Readdirnames(-1)
	if err != nil {
		t.Fatalf("Readdir should not throw error.")
	}
	sort.Strings(names)
	ExpectEqual(t, "a", names[0])
	ExpectEqual(t, "b.txt", names[1])
	ExpectEqual(t, "c.txt", names[2])
}

func TestFileChdir(t *testing.T) {
	mf := NewMockFilesystem()
	mf.MkdirAll("/foo/bar/baz", 0755)
	mf.Create("/foo/bar/baz/foo.txt")
	f, _ := mf.Open("/foo/bar/baz/foo.txt")
	ExpectCwd(t, "/", mf)
	f.Chdir()
	ExpectCwd(t, "/foo/bar/baz", mf)
}

func TestFileChmod(t *testing.T) {
	var (
		f   File
		fi  os.FileInfo
		err error
	)
	mf := NewMockFilesystem()
	mf.Create("foo.txt")
	if f, err = mf.Open("foo.txt"); err != nil {
		t.Fatalf("File expected: %v", err)
	}
	fi = ExpectFile(t, "foo.txt", mf)
	if perm := fi.Mode().Perm(); perm != 0666 {
		t.Fatalf("New file perm %v, expected 0666", perm)
	}
	if err = f.Chmod(0755); err != nil {
		t.Fatalf("Chmod should not return error: %v", err)
	}
	if perm := fi.Mode().Perm(); perm != 0755 {
		t.Fatalf("Perm %v, expected 0755", perm)
	}
}

func TestStat(t *testing.T) {
	var (
		mf  *MockFilesystem
		f   File
		fi  os.FileInfo
		err error
	)
	mf = NewMockFilesystem()
	mf.Mkdir("foo", 0755)
	mf.Create("/foo/foo.txt")
	if f, err = mf.Open("/foo/foo.txt"); err != nil {
		t.Fatalf("File should exist: %v", err)
	}
	if fi, err = f.Stat(); err != nil {
		t.Fatalf("Stat should not throw error: %v", err)
	}
	if fi.Name() != "foo.txt" {
		t.Fatalf("Stat should return accurate file object.")
	}
}

func TestFileInfoPath(t *testing.T) {
	mf := NewMockFilesystem()
	mf.MkdirAll("/foo/bar/baz", 0755)
	fp := "/foo/bar/baz/foo.txt"
	mf.Create(fp)
	fi := ExpectFile(t, fp, mf)
	if fi.path() != fp {
		t.Fatalf("Expected path of %v, got %v", fp, fi.path())
	}
}

func TestWrite(t *testing.T) {
	var (
		fs     Filesystem
		f      File
		err    error
		info   os.FileInfo
		input  string
		output []byte
	)
	fs = NewMockFilesystem()
	f, _ = fs.Create("foo.txt")
	input = "Hello world"
	f.Write([]byte(input))
	f.Close()
	if f, err = fs.Open("foo.txt"); err != nil {
		t.Fatalf("Could not open written file: %v", err)
	}
	if info, err = f.Stat(); err != nil {
		t.Fatalf("Stat should not raise error: %v", err)
	}
	if info.Size() != int64(len(input)) {
		t.Fatalf("File size %v, expected %v", info.Size(), len(input))
	}
	output = make([]byte, info.Size())
	if _, err = f.Read(output); err != nil {
		t.Fatalf("Read should not return error: %v", err)
	}
	if string(output) != input {
		t.Fatalf("Read: %v != expected: %v", string(output), input)
	}
}
