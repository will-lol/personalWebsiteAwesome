package blog

import (
	"log/slog"
	"time"
)

type Blog struct {
	Posts []Post 
}

type Post struct {
	Date time.Time
	Title string
	Content string
}

type Files interface {
	GetAllFiles() (*[]*SimpleFile, error)
}

type SimpleFile struct {
	Bytes []byte
	Name string
}

func NewBlog(fileGetter Files) (*Blog, error) {
	files, err := fileGetter.GetAllFiles()
	if err != nil {
		return nil, err
	}
	for _, file := range *files {
		slog.Default().Info(string(file.Bytes))
	}
	return nil, nil
}
