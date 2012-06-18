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
	Readdir(n int) (fi []os.FileInfo, err error)
	Readdirnames(n int) (names []string, err error)
}

type Filesystem struct {}

func (f *Filesystem) Chdir(dir string) error {
	return os.Chdir(dir)
}

func (f *Filesystem) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (f *Filesystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (f *Filesystem) Remove(name string) error {
	return os.Remove(name)
}

func (f *Filesystem) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (f *Filesystem) Rename(oldname string, newname string) error {
	return os.Rename(oldname, newname)
}

func (f *Filesystem) Create(name string) (file File, err error) {
	return os.Create(name)
}

func (f *Filesystem) Open(name string) (file File, err error) {
	return os.Open(name)
}

func (f *Filesystem) OpenFile(name string, flag int, perm os.FileMode) (file File, err error) {
	return os.OpenFile(name, flag, perm)
}
