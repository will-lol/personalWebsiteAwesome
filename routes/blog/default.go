package blog

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	blogService "github.com/will-lol/personalWebsiteAwesome/services/blog"
)

type BlogHandler struct {
	Files       *blogService.Files
	Log         *slog.Logger
	blogService *blogService.Blog
}

func NewBlogHandler(f *blogService.Files, l *slog.Logger) (*BlogHandler, error) {
	bService, err := blogService.NewBlog(*f)
	if err != nil {
		return nil, err
	}

	return &BlogHandler{
		Files: f,
		Log:   l,
		blogService: bService,
	}, nil
}

func (b BlogHandler) blogHandler(w http.ResponseWriter, r *http.Request) {
	component := blog()
	component.Render(r.Context(), w)
}

func (b BlogHandler) Router(r chi.Router) {
	r.Get("/", b.blogHandler)
}
