package fs

import (
	"os"

	"github.com/will-lol/personalWebsiteAwesome/services/blog"
)

type Fs struct {
	folderName string
}

func NewFs(folder string) *Fs {
	return &Fs{
		folderName: folder,
	}
}

func (f Fs) GetAllFiles() (*[]*blog.SimpleFile, error) {
	files, err := os.ReadDir(f.folderName)
	if err != nil {
		return nil, err
	}

	out := make([]*blog.SimpleFile, len(files), len(files))
	for i, file := range files {
		name := file.Name()
		bytes, err := os.ReadFile(f.folderName + "/" + name)
		if err != nil {
			return nil, err
		}
		out[i] = &blog.SimpleFile{
			Bytes: bytes,
			Name: name,
		}
	}

	return &out, nil
}
