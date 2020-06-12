package files

import (
	"fmt"
	"io/ioutil"

	"github.com/thamizhv/tgnutella/helpers"
	"github.com/thamizhv/tgnutella/models"
)

type FileHandler interface {
	Count() uint32
	Size() uint32
	Get(key string) *models.File
	Exists(key string) bool
	UpdateFileList()
}

type fileHelper struct {
	dir     string
	count   uint32
	size    uint32
	changed bool
	files   map[string]*models.File
}

func NewFileHelper(dir string) FileHandler {
	fileHelper := &fileHelper{
		dir:     dir,
		changed: true,
		files:   make(map[string]*models.File),
	}

	fileHelper.update()

	return fileHelper
}

func (f *fileHelper) Count() uint32 {
	f.update()
	return f.count
}

func (f *fileHelper) Size() uint32 {
	f.update()
	return f.size
}

func (f *fileHelper) Exists(key string) bool {
	if _, ok := f.files[key]; ok {
		return true
	}

	return false
}

func (f *fileHelper) Get(key string) *models.File {
	val, ok := f.files[key]
	if !ok {
		return nil
	}

	return val
}

func (f *fileHelper) UpdateFileList() {
	f.changed = true
	f.update()
}

func (f *fileHelper) update() {
	if !f.changed {
		return
	}

	files, err := ioutil.ReadDir(f.dir)
	if err != nil {
		fmt.Printf("Error in reading current directory %s: %v\n", f.dir, err)
		return
	}

	var size int64

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fi := &models.File{
			Name: file.Name(),
			Size: file.Size(),
		}
		id := helpers.GetHash(file.Name())
		f.files[id] = fi
		size += file.Size()
	}

	f.size = uint32(size)
	f.count = uint32(len(f.files))
	f.changed = false
}
