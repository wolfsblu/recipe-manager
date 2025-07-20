package routing

import (
	"net/http"

	swagger "github.com/swaggest/swgui/v5emb"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/wolfsblu/go-chef/api"
	"github.com/wolfsblu/go-chef/infra/config"
	"github.com/wolfsblu/go-chef/infra/env"
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
	mux.Handle(config.ImagesPathPrefix+"/", http.StripPrefix(config.ImagesPathPrefix+"/", fsHandler))
}

func handleAPI(mux *http.ServeMux, apiServer http.Handler) {
	mux.Handle(config.APIPathPrefix+"/docs/", swagger.New("OpenAPI Docs", config.APIPathPrefix+"/openapi.yml", config.APIPathPrefix+"/docs/"))
	mux.Handle(config.APIPathPrefix+"/", cors(http.StripPrefix(config.APIPathPrefix, apiServer)))
	mux.HandleFunc(config.APIPathPrefix+"/openapi.yml", apiDocs)
}

func handleUploads(mux *http.ServeMux, uploadServer *tusd.Handler) {
	mux.Handle(config.UploadPathPrefix+"/", cors(http.StripPrefix(config.UploadPathPrefix+"/", uploadServer)))
	mux.Handle(config.UploadPathPrefix, cors(http.StripPrefix(config.UploadPathPrefix, uploadServer)))
}
