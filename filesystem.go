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
)

type File interface {
	Chdir() error
	Chmod(mode os.FileMode) error
	Close() error
	Name() string
	Read(b []byte) (n int, err error)
	ReadAt(b []byte, off int64) (n int, err error)
	Readdir(n int) (fi []os.FileInfo, err error)
	Readdirnames(n int) (names []string, err error)
	Stat() (fi os.FileInfo, err error)
	Sync() (err error)
	Seek(offset int64, whence int) (ret int64, err error)
	Truncate(size int64) error
	Write(b []byte) (n int, err error)
	WriteAt(b []byte, off int64) (n int, err error)
	WriteString(s string) (ret int, err error)
}

type Filesystem interface {
	Chdir(dir string) error
	Mkdir(name string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Remove(name string) error
	RemoveAll(path string) error
	Rename(oldname string, newname string) error
	Create(name string) (file File, err error)
	Open(name string) (file File, err error)
	OpenFile(name string, flag int, perm os.FileMode) (file File, err error)
	Stat(name string) (fi os.FileInfo, err error)
}

type RealFilesystem struct {}

func (f *RealFilesystem) Chdir(dir string) error {
	return os.Chdir(dir)
}

func (f *RealFilesystem) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (f *RealFilesystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (f *RealFilesystem) Remove(name string) error {
	return os.Remove(name)
}

func (f *RealFilesystem) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (f *RealFilesystem) Rename(oldname string, newname string) error {
	return os.Rename(oldname, newname)
}

func (f *RealFilesystem) Create(name string) (file File, err error) {
	return os.Create(name)
}

func (f *RealFilesystem) Open(name string) (file File, err error) {
	return os.Open(name)
}

func (f *RealFilesystem) OpenFile(name string, flag int, perm os.FileMode) (file File, err error) {
	return os.OpenFile(name, flag, perm)
}

func (f *RealFilesystem) Stat(name string) (fi os.FileInfo, err error) {
	return os.Stat(name)
}
