package fs

import (
	"os"

	"github.com/will-lol/personalWebsiteAwesome/services/blog"
)

type fs struct {
	folderName string
}

func NewFs(folder string) *fs {
	return &fs{
		folderName: folder,
	}
}

func (f fs) GetAllFiles() (*[]*blog.SimpleFile, error) {
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
