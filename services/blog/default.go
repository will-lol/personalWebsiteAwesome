package blog

import (
	"bytes"
	"errors"
	"strings"
	"sync"
	"time"

	fm "github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

type blog struct {
	Posts *map[string]*Post
}

type Blog interface {
	GetPost(slug string) (*Post, error)
	GetAllPosts() *map[string]*Post
}

type frontmatter struct {
	Date        time.Time `yaml:"date"`
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	Slug        string    `yaml:"slug"`
}

type Post struct {
	Frontmatter frontmatter
	Content     string
}

type Files interface {
	GetAllFiles() (*[]*SimpleFile, error)
}

type SimpleFile struct {
	Bytes []byte
	Name  string
}

func NewBlog(fileGetter Files) (*blog, error) {
	files, err := fileGetter.GetAllFiles()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	postCh := make(chan *Post, len(*files))
	errCh := make(chan error, len(*files))

	for _, file := range *files {
		wg.Add(1)

		go func(file *SimpleFile, errCh chan error, postCh chan *Post) {
			defer wg.Done()

			post, err := parsePost(file)
			if err != nil {
				errCh <- err
			}
			postCh <- post
		}(file, errCh, postCh)
	}

	go func() {
		wg.Wait()
		close(errCh)
		close(postCh)
	}()

	for err := range errCh {
		return nil, err
	}

	posts := make(map[string]*Post, len(*files))
	for post := range postCh {
		posts[post.Frontmatter.Slug] = post
	}

	return &blog{
		Posts: &posts,
	}, nil
}

func parsePost(file *SimpleFile) (*Post, error) {
	var matter frontmatter
	rest, err := fm.Parse(strings.NewReader(string(file.Bytes)), &matter)
	if err != nil {
		return nil, err
	}

	var buff bytes.Buffer

	md := goldmark.New(
		goldmark.WithRendererOptions(html.WithUnsafe()),
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("pygments"),
			),
			extension.Typographer,
		),
	)

	if err := md.Convert(rest, &buff); err != nil {
		return nil, err
	}

	return &Post{
		Frontmatter: matter,
		Content:     buff.String(),
	}, nil
}

func (b blog) GetPost(slug string) (*Post, error) {
	post := (*b.Posts)[slug]
	if post != nil {
		return post, nil
	}
	return nil, errors.New("Post not found")
}

func (b blog) GetAllPosts() *map[string]*Post {
	return b.Posts
}
