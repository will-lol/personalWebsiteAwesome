package index

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"log"
	"context"
	"github.com/will-lol/personal_website/eid"
)

func GetIdFactory(ctx context.Context) (*eid.EidFactory) {
	if factory, ok := ctx.Value("idFactory").(eid.EidFactory); ok {
		return &factory
	}
	log.Fatalln("factory not in context")
	return nil
}

func root(w http.ResponseWriter, r *http.Request) {
	idFactory := eid.New()
	component := index()
	ctx := context.WithValue(r.Context(), "idFactory", idFactory)
	component.Render(ctx, w)
}

func Router(r chi.Router) {
	r.Get("/", root)
}
