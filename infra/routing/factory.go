package routing

import (
	swagger "github.com/swaggest/swgui/v5emb"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/wolfsblu/go-chef/api"
	"github.com/wolfsblu/go-chef/infra/env"
	"net/http"
)

func NewServeMux(server *api.Server, uploadServer *tusd.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	handleFrontend(mux)
	handleImages(mux)
	handleUploads(mux, uploadServer)
	handleAPI(mux, server)
	return mux
}

func handleFrontend(mux *http.ServeMux) {
	mux.HandleFunc("/assets/", assets)
	mux.HandleFunc("/", index)
}

func handleImages(mux *http.ServeMux) {
	imageFs := newFileOnlyFileSystem(env.MustGet("IMAGE_PATH"))
	fsHandler := http.FileServer(imageFs)
	mux.Handle("/images/", http.StripPrefix("/images", fsHandler))
}

func handleAPI(mux *http.ServeMux, apiServer http.Handler) {
	mux.Handle("/api/docs/", swagger.New("OpenAPI Docs", "/api/openapi.yml", "/api/docs/"))
	mux.Handle("/api/", cors(http.StripPrefix("/api", apiServer)))
	mux.HandleFunc("/api/openapi.yml", apiDocs)
}

func handleUploads(mux *http.ServeMux, uploadServer *tusd.Handler) {
	mux.Handle("/api/uploads/", cors(http.StripPrefix("/api/uploads/", uploadServer)))
	mux.Handle("/api/uploads", cors(http.StripPrefix("/api/uploads", uploadServer)))
}
