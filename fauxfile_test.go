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
	if mf.cwd != expected {
		t.Fatalf("Expected cwd of '%v', got '%v'", expected, mf.cwd)
	}
}

func TestChdir(t *testing.T) {
	mf := NewMockFilesystem()
	mf.Mkdir("/foo", 0755)
	ExpectCwd(t, "/", mf)
	mf.Chdir("foo")
	ExpectCwd(t, "/foo", mf)
}

