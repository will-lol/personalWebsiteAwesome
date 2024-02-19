package blog

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	blogService "github.com/will-lol/personalWebsiteAwesome/services/blog"
)

type blogHandler struct {
	FilesDependency *blogService.Files
	Log             *slog.Logger
	blogService     blogService.Blog
}

func NewBlogHandler(f *blogService.Files, l *slog.Logger) (*blogHandler, error) {
	bService, err := blogService.NewBlog(*f)
	if err != nil {
		return nil, err
	}

	return &blogHandler{
		FilesDependency: f,
		Log:             l,
		blogService:     *bService,
	}, nil
}

func (b blogHandler) blogHandler(w http.ResponseWriter, r *http.Request) {
	component := blog(b.blogService)
	component.Render(r.Context(), w)
}

func (b blogHandler) postHandler(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	blogPost, err := b.blogService.GetPost(slug)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("404 Not Found"))
		return
	}

	component := postPage(blogPost)
	component.Render(r.Context(), w)
}

func (b blogHandler) Router(r chi.Router) {
	r.Get("/", b.blogHandler)
	r.Get("/{slug}", b.postHandler)
}

